from django.db import models

from server.user.models import User

# Create your models here.
class UserSession(models.Model):
    token = models.CharField
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    exp = models.BigIntegerField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        db_table = "user_session"