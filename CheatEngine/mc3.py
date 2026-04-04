from flask import Flask, request, jsonify
from pyngrok import ngrok
from transformers import TextStreamer

app = Flask(__name__)

ngrok.set_auth_token("2OGrdSJJOMfZVL4lsN9ES4a7NQX_3mHgdmcwYPK4ZypaRvqca")  # Replace with your actual Ngrok auth token

@app.route("/message", methods=["POST"])
def final_message():
    data = request.json
    question=data['question']
    schema=data['table']
    query=data['query']
    # alpaca_prompt = Copied from above
    FastLanguageModel.for_inference(model) # Enable native 2x faster inference
    inputs = tokenizer(
    [
    alpaca_prompt2.format(question,schema,query)
    ], return_tensors = "pt").to("cuda")


    outputs = model.generate(**inputs, max_new_tokens=200)
    generated_text = tokenizer.decode(outputs[0], skip_special_tokens=True)
    print(generated_text)
    return jsonify(message=exxtraction2(generated_text))
    
@app.route("/", methods=["POST"])
def hello_world():
    data = request.json
    question=data['question']
    schema=data['schema']
    
    print("Received data:", data)  
    FastLanguageModel.for_inference(model) # Enable native 2x faster inference
    inputs = tokenizer(
    [
        alpaca_prompt.format(
            question, # instruction
            schema, # input
            "", # output - leave this blank for generation!
        )
    ], return_tensors = "pt").to("cuda")


    outputs = model.generate(**inputs, max_new_tokens=40)
    generated_text = tokenizer.decode(outputs[0], skip_special_tokens=True)
    print(exxtraction(generated_text))
    return jsonify(message=exxtraction(generated_text))  
if __name__ == "__main__":
    public_url = ngrok.connect(5000, bind_tls=True)  
    print(f" * ngrok tunnel \"{public_url}\" -> \"http://127.0.0.1:5000/\"")

    app.run()
