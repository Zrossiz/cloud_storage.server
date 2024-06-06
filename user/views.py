from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
import bcrypt
import datetime
import dateutil.relativedelta
import uuid
from django.http import HttpResponse
from .models import User
from .serializers import UserSerializer, UserSessionSerializer

class UserView(APIView):

    def post(self, request, *args, **kwargs):
        if 'sign-up' in request.path:
            return self.sign_up(request)
        
        return Response({'message': 'Invalid endpoint'}, status=status.HTTP_404_NOT_FOUND)

    def sign_up(self, request):
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
            session = self._generate_token(serializer.data['id'])

            response =  Response({
                'success': True,
                'data': serializer.data,
            }, status=status.HTTP_200_OK)

            response.set_cookie(
                key='session',
                value=session,
                httponly=True,
                samesite='Strict',
                secure=True
            )

            return response

        except Exception as e:
            print(e)
            return Response({
                'success': False,
                'message': 'Server error',
            }, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

    def login(self, request):
        return HttpResponse('login route')

    def logout(self, request):
        return HttpResponse('logout route')

    def logout_session(self, request):
        return HttpResponse('logout session route')

    def delete(self, request):
        return HttpResponse('delete user route')

    def _generate_token(self, user_id):
        cur_time = datetime.datetime.now()
        next_month = cur_time + dateutil.relativedelta.relativedelta(months=1)
        diff_time = next_month - cur_time
        data = {
            'token': uuid.uuid4(),
            'exp': int(diff_time.total_seconds() * 1000), 
            'user': user_id
        }

        serializer = UserSessionSerializer(data=data)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return data['token']
