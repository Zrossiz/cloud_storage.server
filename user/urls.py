from django.urls import path
from .views import UserView

urlpatterns = [
    path('sign-up/', UserView.sign_up),
    path('login/', UserView.login),
    path('logout/', UserView.logout),
    path('logout/<int:pk>/', UserView.logout_session),
    path('delete/<int:pk>/', UserView.delete)
]