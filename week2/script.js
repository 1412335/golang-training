import http from 'k6/http';

export let options = {
    vus: 10,
    duration: '30s',
};

export default function () {
    var url = 'http://goapp:8080/increase';
    http.get(url);
}