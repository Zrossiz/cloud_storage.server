from rest_framework import generics, permissions
from django.contrib.auth.models import User
from .serializers import RegisterSerializer, UserSerializer, LoginSerializer
from rest_framework.response import Response
from rest_framework.views import APIView
from django.contrib.auth import login, logout

class RegisterView(generics.CreateAPIView):
    queryset = User.objects.all()
    permission_classes = (permissions.AllowAny,)
    serializer_class = RegisterSerializer

class UpdateProfileView(generics.RetrieveUpdateAPIView):
    queryset = User.objects.all()
    permission_classes = [permissions.IsAuthenticated]
    serializer_class = UserSerializer

    def get_object(self):
        return self.request.user

class LoginView(APIView):
    permission_classes = (permissions.AllowAny,)

    def post(self, request, *args, **kwargs):
        serializer = LoginSerializer(data=request.data)
        if serializer.is_valid():
            user = serializer.validated_data['user']
            login(request, user)
            return Response({'status': 'login successful'})
