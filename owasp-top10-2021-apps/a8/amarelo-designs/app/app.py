# coding: utf-8

from flask import Flask, request, make_response, render_template, redirect, flash
import jwt
import os
from datetime import datetime, timedelta, timezone

app = Flask(__name__)


admin = os.getenv("USERNAME")
admin_pass = os.getenv("PASSWORD")
secret = os.getenv("SECRET")

@app.route("/")
def ola():
    return render_template('index.html')

@app.route("/admin", methods=['GET','POST'])
def login():
    if request.method == 'POST':
        username = request.values.get('username')
        password = request.values.get('password')
    
        if username == admin and password == admin_pass:

            payload = {
                "admin" : True,
                "exp": datetime.now(timezone.utc) + timedelta(minutes=10)
            }
            token = jwt.encode(payload, secret, algorithm="HS256")
            
            resp = make_response(redirect("/user"))
            resp.set_cookie("sessionId", token, httponly=True, secure=False)
            
            return resp

        else:
            return redirect("/admin")
        

    else:
        return render_template('admin.html')

@app.route("/user", methods=['GET'])
def userInfo():
    token = request.cookies.get("sessionId")
    if not token:
        return "Não Autorizado!"
    
    try: 
        data = jwt.decode(token, secret, algorithms=["HS256"])
        return render_template('user.html', user=data)

    except jwt.ExpiredSignatureError:
        return "Sessão expirada, faça login novamente."
    
    except jwt.InvalidTokenError:
        return "Token inválido!"
    
if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
