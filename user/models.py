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