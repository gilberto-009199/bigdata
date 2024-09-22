var api = {
    socket:(containerId) => new WebSocket("ws://localhost:8080/ws?containerId="+containerId),
    start: (containerId) => fetch("/api/start/?containerId="+containerId).then(response => response.json()),
    stop: (containerId) => fetch("/api/stop/?containerId="+containerId).then(response => response.json()),
    duplicar:(containerId) => fetch("/api/duplicar/?containerId="+containerId).then(response => response.json()),
    fetchUpload: (formData) => fetch('/upload', {
        method: 'POST',
        body: formData,
    }),
    fecthApi: () => fetch('/api').then(response => response.json())
}