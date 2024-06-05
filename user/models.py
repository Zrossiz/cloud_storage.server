from django.db import models

# Create your models here.
class User(models.Model):
    username = models.CharField()
    email = models.CharField(blank=True, null=True)
    password = models.CharField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "user"

class UserSession(models.Model):
    token = models.CharField
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    exp = models.BigIntegerField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "user_session"