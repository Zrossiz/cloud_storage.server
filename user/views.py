from django.http import HttpResponse
from rest_framework.response import Response
from rest_framework import status
from rest_framework.decorators import api_view
import bcrypt

from .serializers import UserSerializer
from .models import User

# Create your views here.
class UserView():

    @api_view(['POST'])
    def sign_up(request):
        try:
            body = request.data
            existUser = User.objects.filter(username=body['username']).exists()
            
            if existUser:
                return Response({
                'success': False,
                'message': 'username or email already taken',
            }, status=status.HTTP_400_BAD_REQUEST)


            password = request.data["password"].encode('utf-8')
            hashed_password = bcrypt.hashpw(password, bcrypt.gensalt())
            request.data['password'] = hashed_password.decode('utf-8')
            serializer = UserSerializer(data=request.data)
            serializer.is_valid(raise_exception=True)
            serializer.save()

            return Response({
                'success': True,
                'data': serializer.data,
            }, status=status.HTTP_200_OK)

        except Exception as e:
            print(e)
            return Response({
                'success': False,
                'message': 'Server error',
            }, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    
    def login(request):
        return HttpResponse('login route')
    
    def logout(request):
        return HttpResponse('logout route')
    
    def logout_session(request):
        return HttpResponse('logout session route')
    
    def delete(request):
        return HttpResponse('delete user route')