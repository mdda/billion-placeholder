## Viewing the Presentation

This presentation makes use of the excellent [```reveal.js``` presentation tool](https://github.com/hakimel/reveal.js).

To view the presentation locally requires that you unpack ```reveal.js``` here appropriately 
(it doesn't require installation on your machine, just unwrapping): 

```
# pwd == docs/2015-01-15_Presentation-PyDataSG/
wget https://github.com/hakimel/reveal.js/archive/2.6.2.tar.gz
tar -xzf 2.6.2.tar.gz 
```

Open the presentation in Firefox or Chrome (or another modern browser) using the path given by : 

```
echo `pwd`/reveal.js-2.6.2/presentation.html
```

### nginx configuration 

There's no need to do this, unless you *really* want to host the 
presentation on a [publicly visible server](RedCatLabs.com/2015-01-15_Presentation-PySG/):

```
nginx.conf ::
        location ~ ^/2015-01-15_Presentation-PyDataSG/ {
                root    ...full-path-to-repo.../docs/2015-01-15_Presentation-PyDataSG/reveal.js-2.6.2/;
                rewrite ^/2015-01-15_Presentation-PyDataSG/$ /presentation.html break;
                rewrite ^/2015-01-15_Presentation-PyDataSG/(.+)$ /$1 break;
        }
```
