from django.db import models

# Create your models here.

class Insta(models.Model):
    url = models.URLField()
    date = models.DateTimeField()
    server = models.CharField(max_length=100)
    
    def __str__(self):
        return self.server