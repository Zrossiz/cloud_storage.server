from django.urls import path
from .views import UserView

urlpatterns = [
    path('sign-up/', UserView.as_view()),
    path('login/', UserView.as_view()),
    path('logout/', UserView.as_view()),
    path('logout/<int:pk>/', UserView.as_view()),
    path('delete/<int:pk>/', UserView.as_view())
]