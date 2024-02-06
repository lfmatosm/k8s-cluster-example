# post-app
Static webpage to manage image posts served through Apache. For testing outside the cluster, you need to run the [`backend`](../../back/post-service) application too.

## Building, tagging and pushing
```sh
docker build -t post-app .
docker tag post-app <my_user>/post-app:latest
docker login -U <my_user>
docker push <my_user>/post-app:latest
```

## Resources

- https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Client-side_web_APIs/Fetching_data
- https://kubernetes.io/docs/tasks/access-application-cluster/connecting-frontend-backend/#creating-the-frontend
- https://httpd.apache.org/docs/2.4/howto/reverse_proxy.html#simple
- https://httpd.apache.org/docs/2.4/mod/mod_proxy.html
