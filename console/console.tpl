apiVersion: apps/v1
kind: Deployment
metadata:
  name: console
  labels:
    app: console
spec:
  replicas: 1
  selector:
    matchLabels:
      app: console
  template:
    metadata:
      labels:
        app: console
    spec:
      containers:
        - env:
          - name: CONSOLE_GOOGLE_KEY
            valueFrom:
              secretKeyRef:
                key: console-google-key
                name: console-secrets
          - name: CONSOLE_GOOGLE_SECRET
            valueFrom:
              secretKeyRef:
                key: console-google-secret
                name: console-secrets
          - name: CONSOLE_DATABASE_CONN
            valueFrom:
              secretKeyRef:
                key: console-database-conn
                name: console-secrets
          - name: CONSOLE_CALLBACK_URL
            valueFrom:
              secretKeyRef:
                key: console-callback-url
                name: console-secrets
          - name: CONSOLE_SESSION_KEY
            valueFrom:
              secretKeyRef:
                key: console-session-key
                name: console-secrets
          name: console
          image: docker-registry.curiosityworks.org/curiosinauts/console:__tag__
          ports:
            - containerPort: 3000
      dnsPolicy: Default

---
apiVersion: v1
kind: Service
metadata:
  name: console
spec:
  selector:
    app: console
  ports:
    - protocol: TCP
      port: 80
      targetPort: 3000

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: console
  annotations:
    kubernetes.io/ingress.class: "traefik"
spec:
  rules:
    - host: console.curiosityworks.org
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: console
                port:
                  number: 80