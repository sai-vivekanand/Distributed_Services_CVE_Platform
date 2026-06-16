import os
from . import create_json_logger
from langchain.chains import RetrievalQA
from langchain.prompts import PromptTemplate
from langchain_pinecone import PineconeVectorStore
from pinecone import Pinecone
from langchain.embeddings import HuggingFaceHubEmbeddings
from langchain_community.llms.ollama import Ollama
from langchain.callbacks.manager import CallbackManager
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler

# Initialize logger
logger = create_json_logger()

index_name = os.getenv("PINECONE_INDEX_NAME")
namespace = os.getenv("PINECONE_NAMESPACE")
pinecone_api_key = os.getenv("PINECONE_API_KEY")
hface_api_token = os.getenv("HFACE_API_TOKEN")
ollama_url = os.getenv("OLLAMA_URL")
model_name = os.getenv("MODEL_NAME")
# print(os.getenv("GROQ_API_KEY"))
# print(os.getenv("PINECONE_API_KEY"))
# print(os.getenv("HFACE_API_TOKEN"))
pinecone = Pinecone(api_key=pinecone_api_key)

# Set up the embedding model to use via Hugging Face Inference API
embeddings = HuggingFaceHubEmbeddings(
    repo_id="sentence-transformers/all-MiniLM-L6-v2",
    huggingfacehub_api_token=hface_api_token
)

# Create an index if it does not exist
existing_indexes = [ index.name for index in pinecone.list_indexes() ]
index = pinecone.Index(index_name)

db = PineconeVectorStore(
    index=index,
    namespace=namespace,
    embedding=embeddings
)
retriever = db.as_retriever()
retriever.search_kwargs = {"k": 5} 

model = Ollama(
    base_url=ollama_url, 
    model = model_name,
    temperature = 0,
    callback_manager = CallbackManager([StreamingStdOutCallbackHandler()])
)

template = """Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.
{context}
Question: {question}
Helpful Answer:"""
QA_CHAIN_PROMPT = PromptTemplate.from_template(template)

# Create the RetrievalQA chain
qa_chain = RetrievalQA.from_chain_type(
    llm=model,
    chain_type="stuff",
    retriever=retriever,
    chain_type_kwargs={"prompt": QA_CHAIN_PROMPT}
)

def generate_response(prompt):
    logger.info(f"Generating response for prompt: {prompt}")
    try:
        # for chunk in qa_chain.stream(prompt):
        #     if not isinstance(chunk, dict):
        #         response += str(chunk)
        response = qa_chain.invoke(prompt)
        logger.info(f"Generated response: {response}")
    except Exception as e:
        logger.error(f"Error generating response: {e}")
        response = f"Error: {e}"
    return response