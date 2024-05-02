from __future__ import annotations
from typing import List
from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.orm import Mapped
from datetime import datetime
import models
from sqlalchemy.orm import relationship
from sqlalchemy import Column
from sqlalchemy import Table
from sqlalchemy import ForeignKey


TeamMembership = Table(
    "team_memberships",
    models.Base.metadata,
    Column("user_id", ForeignKey("users.id")),
    Column("team_id", ForeignKey("teams.id")),
    # Column("created_at",DateTime(), default=datetime.now, onupdate=datetime.now)
)


class User(models.Base):


    __tablename__ = 'users'

    id = Column(Integer(), primary_key=True)
    email = Column(String(256), nullable=False, unique=True)
    first_name = Column(String(256), nullable=False)
    last_name = Column(String(256), nullable=False)
    picture_url = Column(String(512), nullable=False)
    created_at = Column(DateTime(), default=datetime.now)
    updated_at = Column(DateTime(), default=datetime.now, onupdate=datetime.now)

    teams: Mapped[List[Team]] = relationship(secondary=TeamMembership)


class Team(models.Base):
    __tablename__ = 'teams'

    id = Column(Integer(), primary_key=True)
    name = Column(String(256), nullable=False)
    created_at = Column(DateTime(), default=datetime.now)
    updated_at = Column(DateTime(), default=datetime.now, onupdate=datetime.now)

    members: Mapped[List[User]] = relationship(secondary=TeamMembership)


