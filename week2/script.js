import http from 'k6/http';

export default function () {
    var url = 'http://goapp:8080/increase';
    http.get(url);
}