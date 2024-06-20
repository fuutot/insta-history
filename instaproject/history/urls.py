from django.urls import path

from . import views
#from .views import read_history

urlpatterns = [
    path('', views.index_view, name='index'),
    path('list/', views.ListInstaView.as_view(), name='list-insta'),
    path('<int:pk>/detail/',views.DetailInstaView.as_view(), name='detail-insta'),
    path('reload/', views.read_history_view, name='reload-insta'),
    
]