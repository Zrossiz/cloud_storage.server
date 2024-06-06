from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
import bcrypt
import datetime
import dateutil.relativedelta
import uuid
from django.http import HttpResponse
from .models import User, UserSession
from .serializers import UserSerializer, UserSessionSerializer

class UserView(APIView):

    def post(self, request, *args, **kwargs):

        if 'sign-up' in request.path:
            return self.sign_up(request)
        if 'login' in request.path:
            return self.login(request)
        
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
            session = self._get_token(serializer.data['id'])

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
        try:
            data = request.data
            username = data['username']
            password = data['password'].encode('utf-8')

            try:
                user = User.objects.get(username=username)
            except User.DoesNotExist:
                return Response({
                    'success': False,
                    'message': 'username or password invalid',
                }, status=status.HTTP_400_BAD_REQUEST)
            
            is_matched_passwords = bcrypt.checkpw(password, user.password.encode('utf-8'))

            if not is_matched_passwords:
                return Response({
                    'success': False,
                    'message': 'username or password invalid',
                }, status=status.HTTP_400_BAD_REQUEST)
            

            serializer = UserSerializer(user)
            session = self._get_token(serializer.data['id'])

            if not session:
                 return Response({
                    'success': False,
                    'message': 'unauthorized',
                }, status=status.HTTP_403_FORBIDDEN)

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

    def logout(self, request):
        return HttpResponse('logout route')

    def logout_session(self, request):
        return HttpResponse('logout session route')

    def delete(self, request):
        return HttpResponse('delete user route')

    def _get_token(self, user_id):
        exist_session = UserSession.objects.filter(user_id=user_id).first()

        if exist_session:
            exp_datetime = datetime.datetime.fromtimestamp(exist_session.exp / 1000.0, datetime.timezone.utc)

            if exp_datetime > datetime.datetime.now(datetime.timezone.utc):
                return exist_session.token
            else:
                exist_session.delete()
                return None

        cur_time = datetime.datetime.now()
        next_month = cur_time + dateutil.relativedelta.relativedelta(months=1)

        data = {
            'token': uuid.uuid4(),
            'exp': int(next_month.timestamp() * 1000), 
            'user': user_id
        }

        serializer = UserSessionSerializer(data=data)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return data['token']
