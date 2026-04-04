from flask import Flask, Response
from pyngrok import ngrok


app = Flask(__name__)

ngrok.set_auth_token("2OGrdSJJOMfZVL4lsN9ES4a7NQX_3mHgdmcwYPK4ZypaRvqca") 

@app.route("/file")
def get_file_content():
    with open("paymentmeta.txt", "r") as f:
        content = f.read()

    return Response(content, mimetype="text/plain")


@app.route("/loan-action/payment/<lan>/metadata", methods=["POST"])
def get_file_content2(lan):
    with open("paymentmeta.txt", "r") as f:
        content = f.read()

    return Response(content, mimetype="text/plain")


@app.route("/loan-action/payment/<lan>/create", methods=["POST"])
def getpamentResponse(lan):
    with open("payment.txt", "r") as f:
        content = f.read()

    return Response(content, mimetype="text/plain")


@app.route("/loan-action/tranche-disbursal/<lan>/create", methods=["POST"])
def getTrancheResponse(lan):
    with open("payment.txt", "r") as f:
        content = f.read()

    return Response(content, mimetype="text/plain")
if __name__ == "__main__":

    # public_url = ngrok.connect(5000, bind_tls=True)  
    # print(f" * ngrok tunnel \"{public_url}\" -> \"http://127.0.0.1:5000/\"")
    app.run(debug=True)
