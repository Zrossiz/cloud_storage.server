from django.http import HttpResponse

# Create your views here.
class UserView():

    def sign_up(request):
        return HttpResponse('sign up route')
    
    def login(request):
        return HttpResponse('login route')
    
    def logout(request):
        return HttpResponse('logout route')
    
    def logout_session(request):
        return HttpResponse('logout session route')
    
    def delete(request):
        return HttpResponse('delete user route')