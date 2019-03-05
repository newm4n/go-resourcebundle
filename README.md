# go-resourcebundle

The java's ResourceBundle implementation in Go.
ResourceBundle is an approach on `i18n` for data text so can it be represented in multiple language.

## How to get

```text
$ go get github.com/newm4n/go-resourcebundle
```

## How it works

Consider the following bundles :

```text
# file : default.properties
greeting=Hello dear customer, welcome to petstore.com
purchase.confirm=Are you sure to purchase {{item}} for {{currency}} {{amount}} ?
purchase.thanks=Thank you for your purchase at petstore.com 
purchase.ok=Ok
purchase.cancel=Cancel

``` 
and
```text
# file : Id_id.properties
greeting=Halo pelangan, selamat datang di petstore.com
purchase.confirm=Apakah anda yakin untuk membeli {{item}} seharga {{currency}} {{amount}} ?
purchase.thanks=Terima kasih untuk pembelian anda di petstore.com 
purchase.cancel=Batal

``` 
and
```text
# file : De_de.properties
greeting=Hallo lieber Kunde, willkommen bei petstore.com
purchase.confirm=Sind Sie sicher, {{item}} für {{currency}} {{amount}} zu kaufen?
purchase.thanks=Vielen Dank für Ihren Einkauf bei petstore.com
``` 

Now you have a dynamic page 

```html
...
<div class="greeting"><% bundle.get("purchase.confirm") %></div>
<input type="button" value="<% bundle.get("purchase.confirm") %>" onClick="confirm()">
<input type="button" value="<% bundle.get("purchase.cancel") %>" onClick="cancel()">
...
```
As you might get the idea, depends on the current user's prefered language, the dynamic page should rendered
according to that (prefered language). `default.properties` represent sets of text belong to any language IF the selected bundle 
dont have the `key` in them.

In the example above `purchase.cancel` are not available if the user's prefered language is `De.de (German)`, the 
the resource bundle should return the value from `default` which is `Cancel`. 

## How to use

TBW