from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
import bcrypt
import uuid
from django.core.cache import cache
from .models import User
from .serializers import UserSerializer

class UserView(APIView):

    def post(self, request, *args, **kwargs):
        if 'sign-up' in request.path:
            return self.sign_up(request)
        if 'login' in request.path:
            return self.login(request)

    def sign_up(self, request):
        try:
            body = request.data
            user_agent = request.META.get('HTTP_USER_AGENT', '')
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
            session = self._get_token(serializer.data['id'], user_agent)

            response = Response({
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
            user_agent = request.META.get('HTTP_USER_AGENT', '')

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
            session = self._get_token(serializer.data['id'], user_agent)

            if not session:
                except_response = Response({
                    'success': False,
                    'message': 'unauthorized',
                }, status=status.HTTP_403_FORBIDDEN)

                except_response.delete_cookie('session')

                return except_response

            response = Response({
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

    def _get_token(self, user_id, user_agent):
        exist_session = cache.get(user_id)
        if exist_session:
            return exist_session

        token = str(uuid.uuid4())
        cache.set(user_id, token, timeout=60*60*24*2)
        return token
