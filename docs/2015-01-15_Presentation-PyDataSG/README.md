## Viewing the Presentation

This presentation makes use of the excellent [```reveal.js``` presentation tool](https://github.com/hakimel/reveal.js).

To view the presentation locally requires that you unpack ```reveal.js``` here appropriately 
(it doesn't require installation on your machine, just unwrapping): 

{% highlight bash %}
# pwd == docs/2015-01-15_Presentation-PySG/
wget https://github.com/hakimel/reveal.js/archive/2.6.2.tar.gz
tar -xzf 2.6.2.tar.gz 
{% endhighlight %}

Open the presentation in Firefox or Chrome (or another modern browser) using the path given by : 

{% highlight bash %}
echo `pwd`/reveal.js-2.6.2/presentation.html
{% endhighlight %}



### nginx configuration 

There's no need to do this, unless you *really* want to host the 
presentation on a [publicly visible server](RedCatLabs.com/2015-01-15_Presentation-PySG/):

{% highlight bash %}
nginx.conf ::
        location ~ ^/2015-01-15_Presentation-PyDataSG/ {
                root    ...full-path-to-repo.../docs/2015-01-15_Presentation-PyDataSG/reveal.js-2.6.2/;
                rewrite ^/2015-01-15_Presentation-PyDataSG/$ /presentation.html break;
                rewrite ^/2015-01-15_Presentation-PyDataSG/(.+)$ /$1 break;
        }
{% endhighlight %}
