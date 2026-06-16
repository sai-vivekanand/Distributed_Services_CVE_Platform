import requests
import streamlit as st
import os

model_name = os.getenv("MODEL_NAME")

st.title("CVE-LLM")

if "messages" not in st.session_state:
    st.session_state.messages = []

for message in st.session_state.messages:
    with st.chat_message(message["role"]):
        st.markdown(message["content"])

if prompt := st.chat_input("What is up?"):
    st.session_state.messages.append({"role": "user", "content": prompt})
    with st.chat_message("user"):
        st.markdown(prompt)

    answer=""
    with st.chat_message("assistant") and st.spinner("Generating..."):
        response = requests.post(
        "http://localhost:5000/generate",
        json={
            "model": model_name,
            "prompt": prompt,
            "stream": False
        }
        )
        response_data = response.json()
        answer = response_data.get("response",{}).get("result", "No response received.")

    st.session_state.messages.append({"role": "assistant", "content": answer})
    st.markdown(answer)