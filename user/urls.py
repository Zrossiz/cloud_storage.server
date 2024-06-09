from django.urls import path
from .views import RegisterView, UpdateProfileView

urlpatterns = [
    path('sign-up/', RegisterView.as_view(), name='register'),
    path('profile/', UpdateProfileView.as_view(), name='profile'),
]
