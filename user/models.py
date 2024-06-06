from django.db import models

# Create your models here.
class User(models.Model):
    username = models.CharField(unique=True)
    email = models.CharField(blank=True, null=True, unique=True)
    password = models.CharField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "user"

class UserSession(models.Model):
    token = models.CharField(max_length=36)
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    exp = models.BigIntegerField()
    user_agent = models.CharField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "user_session"