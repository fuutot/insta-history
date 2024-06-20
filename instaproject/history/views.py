from django.shortcuts import render
from django.views.generic import ListView, DetailView
from django.http import HttpResponse
from .models import Insta
# from django.core.paginator import Paginator
import sqlite3
from contextlib import closing
from datetime import datetime, timezone, timedelta

# from .consts import ITEM_PER_PAGE

# Create your views here.

JST = timezone(timedelta(hours=+9), 'JST')

class ListInstaView(ListView):
    template_name = 'insta/insta_list.html'
    model = Insta
    def get_queryset(self, **kwargs):
        queryset = super().get_queryset(**kwargs)
        query = self.request.GET
        q1 = query.get('q1')
        q2 = query.get('q2')
        if q1 and q2:
            queryset = queryset.filter(date__range=[q1, q2])
        elif q1:
            queryset = queryset.filter(date__gte=q1)
        elif q2:
            queryset = queryset.filter(date__lte=q2)
        else:
            queryset = queryset.all()
        return queryset.order_by('-date')
    
class DetailInstaView(DetailView):
    template_name = 'insta/insta_detail.html'
    #template_name = 'insta/detail.html'
    model = Insta    
    
def index_view(request):
    object_list = Insta.objects.order_by('-id')
    #paginator = Paginator(object_list, ITEM_PER_PAGE)
    #page_number = request.GET.get('page', 1)
    #page_obj = paginator.page(page_number)
    return render(
        request,
        'insta/index.html',
        {'object_list': object_list, 
         #'page_obj': page_obj
         },
        )

db = "/mnt/c/Users/tomoy/AppData/Local/Google/Chrome/User Data/Default/History"
def read_history_view(request):
    Insta.objects.all().delete()
    with closing(sqlite3.connect(db)) as conn:
        c = conn.cursor()
        select_sql = "select visits.id, urls.url, urls.title, visits.visit_time,visits.from_visit from visits inner join urls on visits.url = urls.id"
        for row in c.execute(select_sql):
            timestamp = datetime.fromtimestamp((row[3]-11644473600000000)/1000000, JST)
            if "instagram.com/p/" in row[1]:
                Insta.objects.create(url = row[1], date = timestamp, server = row[2])
    return render(request, 'insta/reload.html')