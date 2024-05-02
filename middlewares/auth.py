from functools import wraps
from flask import request
from sqids import Sqids
from models import db
from authlib.jose import jwt
from models.user import User

sqids = Sqids(min_length=10)
header = {'alg': 'RS256'}
key_file = open('keys/jwt-key')
private_key = key_file.read()
key_file = open('keys/jwt-key.pub')
public_key = key_file.read()


def token_required(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        token = None
        if "Authorization" in request.headers:
            token = request.headers["Authorization"].split(" ")[1]
        if not token:
            return {
                "message": "Request is unauthorized!",
            }, 401

        claims = jwt.decode(token, public_key)
        decoded_user_id = sqids.decode(claims['sub'])[0]
        current_user = db.session.query(User).filter(User.id == decoded_user_id).first()
        if current_user is None:
            return {
                "message": "Request is unauthorized!",
            }, 401

        return f(current_user, *args, **kwargs)

    return decorated
