var api = {
    fecthApi: () => fetch('/api').then(response => response.json())
}