from .models import User, UserSession
from rest_framework import serializers

class UserSerializer(serializers.ModelSerializer):
    
    class Meta:
        model=User
        fields="__all__"

class UserSessionSerializer(serializers.ModelSerializer):

    class Meta:
        model=UserSession
        fields="__all__"