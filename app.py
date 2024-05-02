from flask import Flask
from flask import jsonify
from flask_cors import CORS
from flask import request
import requests
import json
from authlib.jose import jwt
from sqids import Sqids
from flask import abort
import tomllib
import middlewares.auth
import models.team
from models import db
from models.user import User, Team

app = Flask(__name__)

app.config.from_file("app-config.example.toml", load=tomllib.load, text=False)
# You can override config here by adding a new line
app.config.from_file("app-config.local.toml", load=tomllib.load, text=False)

CORS(app)

print()

app.config['SQLALCHEMY_DATABASE_URI'] = ('postgresql+psycopg2://{user}:{pw}@{host}:{port}/{db}'.
                                         format(user=app.config['DATABASE']['USERNAME'],
                                                pw=app.config['DATABASE']['PASSWORD'],
                                                host=app.config['DATABASE']['HOST'],
                                                port=app.config['DATABASE']['PORT'],
                                                db=app.config['DATABASE']['DB_NAME']))
db.init_app(app)

with app.app_context():
    db.create_all()

sqids = Sqids(min_length=10)
header = {'alg': 'RS256'}
key_file = open('keys/jwt-key')
private_key = key_file.read()
key_file = open('keys/jwt-key.pub')
public_key = key_file.read()


@app.route('/v1/teams', methods=['GET'])
@middlewares.auth.token_required
def team_list(user):
    r = []
    for team in user.teams:
        r.append({"name": team.name, "id": team.id})
    return jsonify(r)


@app.route('/v1/login/google', methods=['POST'])
def login_google():
    google_code = request.json['code']

    url = "https://oauth2.googleapis.com/token"

    payload = json.dumps({
        "code": google_code,
        "client_id": app.config['GOOGLE_OAUTH_CLIENT_ID'],
        "client_secret": app.config['GOOGLE_OAUTH_CLIENT_SECRET'],
        "redirect_uri": app.config['GOOGLE_OAUTH_REDIRECT_URI'],
        "grant_type": "authorization_code"
    })
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.post(url, headers=headers, data=payload)
    google_access_token = response.json()['access_token']

    headers = {
        'Authorization': 'Bearer ' + google_access_token
    }
    response = requests.get("https://www.googleapis.com/oauth2/v3/userinfo", headers=headers)
    json_response = response.json()

    u = db.session.query(User).filter(User.email == json_response['email']).first()
    if u is None:
        u = User(email=json_response['email'],
                 first_name=json_response['given_name'],
                 last_name=json_response['family_name'],
                 picture_url=json_response['picture'])
        db.session.add(u)
        db.session.commit()

    encoded_id = sqids.encode([u.id])
    p = {'iss': 'shirinibede', 'sub': encoded_id}

    access_token = jwt.encode(header, p, private_key)

    return jsonify({
        'access_token': access_token.decode(),
        'user': {
            'first_name': u.first_name,
            'last_name': u.last_name,
            'email': u.email,
            'picture_url': u.picture_url,
        }
    })


@app.route('/v1/profile', methods=['GET'])
@middlewares.auth.token_required
def profile(user):
    return jsonify({
        'first_name': user.first_name,
        'last_name': user.last_name,
        'email': user.email,
        'picture_url': user.picture_url,
    })


@app.route('/v1/teams', methods=['POST'])
@middlewares.auth.token_required
def create_team(user):
    name = request.json['name']
    team = models.user.Team(
        name=name,
        members=[user]
    )
    db.session.add(team)
    db.session.commit()

    return jsonify({
        'id': team.id,
        'name': team.name,
        'created_at': team.created_at,
        'updated_at': team.updated_at,
    })


@app.route('/v1/teams/<id>', methods=['GET'])
@middlewares.auth.token_required
def get_team_by_id(user, id):
    team = db.session.query(Team).filter(Team.id == id).first()
    if not team:
        abort(404, 'Team not found')

    return jsonify({
        'id': team.id,
        'name': team.name,
        'created_at': team.created_at,
        'updated_at': team.updated_at,
    })


if __name__ == '__main__':
    app.run()
