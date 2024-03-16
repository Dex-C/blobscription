import urllib.parse
import json
import os
import requests  # Import the requests library
from pymongo import MongoClient  # Import the MongoClient class from pymongo

# Assuming '.credential' is a JSON file, it should be opened with the open() function, then parsed with json.load()
with open('.credential') as f:
    credential = json.load(f)

def format_blobscan_query_url(base_url, json_data):
    json_str = json.dumps(json_data)
    encoded_json_str = urllib.parse.quote(json_str)
    full_url = base_url + encoded_json_str
    return full_url
     
def get_blob_details(commit_hash):
    base_url = "https://sepolia.blobscan.com/api/trpc/blob.getByBlobIdFull?input="
    query_json = {
                "json": {
                    "id": commit_hash
                }
            }
    url = format_blobscan_query_url(base_url, query_json)
    response = requests.get(url)
    if response.status_code == 200:
        data = response.json()  # Assuming the API returns JSON data
        encoded_json_string=data['result']['data']['json']['data']
        encoded_json_string = encoded_json_string[2:]  # Remove the "0x" prefix
        binary_data = bytes.fromhex(encoded_json_string)
        json_string = binary_data.decode('utf-8')
        json_data = json.loads(json_string)
        return json_data

def start_index():
    base_url = "https://sepolia.blobscan.com/api/trpc/tx.getByAddress?input="
    for page in range(1, 1001):  # Assuming pages start from 1
        query_json = {
            "json": {
                # "address": "0x0000000000000000000000000000000000c0ffee",
                "address": "0xff00000000000000000000000000000000009252",
                "p": page,  # Corrected to iterate pages
                "ps": 25
            }
        }
        url = format_blobscan_query_url(base_url, query_json)
        response = requests.get(url)  # Make a GET request to the URL
        if response.status_code == 200:
            data = response.json()  # Parse the JSON response
            for item in data['result']['data']['json']['transactions']:
                print(item)
                block_no=item['block']['number']
                tx_hash=item['hash']
                for blob in item['blobs']:
                    commitment_hash = blob['blob']['commitment']
                    BS20_json=get_blob_details(commitment_hash)  # Print each item in the JSON
                    write_to_mongodb(block_no,commitment_hash,BS20_json)

def start_index_dummy():
    dummy_data_dir = './dummydata'
    for dummy_blob_filename in os.listdir(dummy_data_dir):
        with open(os.path.join(dummy_data_dir, dummy_blob_filename)) as f:
            BS20_json = json.load(f)
            #TODO:generate a dummy block number of 
            #TODO:generate a dummy KZG_commit of length equal to "0x87cec3b5c334e3f89f9594c175cd98269fc16da2179c0a4e8347e071942857558fc17f01cc57f11118468d0b35deedf5"
            write_to_mongodb(BS20_json)

def track_balance():
    client = MongoClient(credential['mongo_uri'])  # Use your MongoDB URI
    db = 'BlobScriptions'  # Specify your database name
    collection = 'BlobScriptions'  # Specify your collection name
    #TODO:retrieve current blocknumber from infura
    #TODO:retrieve all records from the database
    #TODO:if the retrieved record is none, return a defaultdict of defaultdict of integer

def write_to_mongodb(block_no,KZG_Hash,BS20_json):
    # Assuming 'credentials' contains MongoDB connection info
    client = MongoClient(credential['mongo_uri'])  # Use your MongoDB URI
    db = 'BlobScriptions'  # Specify your database name
    collection = 'BlobScriptions'  # Specify your collection name
    #TODO: insert if only KZG commit is not duplicate
    # collection.insert_one(blob_data)  # Insert the JSON data into MongoDB