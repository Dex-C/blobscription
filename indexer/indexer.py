import urllib.parse
import json
import os
import requests  # Import the requests library
from pymongo import MongoClient  # Import the MongoClient class from pymongo
import random
import hashlib
from collections import defaultdict

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
                from_address=item['fromId']
                for blob in item['blobs']:
                    commitment_hash = blob['blob']['commitment']
                    BS20_json=get_blob_details(commitment_hash)  # Print each item in the JSON
                    write_to_mongodb(block_no,commitment_hash,BS20_json,fron_address)

def get_current_block_no():
    infura_url = credential['INFURA']
    response = requests.post(infura_url, json={"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1})
    current_block_number = int(response.json().get('result', '0x0'), 16)
    return current_block_number

def start_index_dummy():
    dummy_data_dir = './dummydata'
    sorted_filenames = sorted(os.listdir(dummy_data_dir), key=lambda x: int(x.replace('.json', '')))
    for dummy_blob_filename in sorted_filenames:
        print(dummy_blob_filename)
        with open(os.path.join(dummy_data_dir, dummy_blob_filename)) as f:
            BS20_json = json.load(f)
            block_no=get_current_block_no()-100+int(dummy_blob_filename.replace('.json',''))
            from_address = credential['EOA']

            data_to_hash = str(block_no) + json.dumps(BS20_json, sort_keys=True)
            # Hashing the combined string
            hash_object = hashlib.sha256(data_to_hash.encode('utf-8'))
            hex_dig = hash_object.hexdigest()
            dummy_KZG_commit = "0x" + hex_dig +"abcdef1234567890abcdef1234567890" # This will be 64 characters long, as SHA-256 produces a 64-character hex digest

            write_to_mongodb(block_no,dummy_KZG_commit,BS20_json,from_address)

def track_balance():
    client = MongoClient(credential['mongo_uri'])  # Use your MongoDB URI
    db = client.BlobScriptions  # Access the database using the client
    collection = db.BlobScriptions  # Access the collection from the database
    current_block_no = get_current_block_no()
    last_4096_epoch = 32 * 4096
    lower_bound_block_no = current_block_no - last_4096_epoch

    # Retrieve records where block_no is greater than current_block_no - last_4096_epoch
    # Sort the records first by block_no, then by KZG_hash, both ascending

    records = list(collection.find({"block_no": {"$gt": lower_bound_block_no}})
                               .sort([("block_no", 1), ("KZG_Hash", 1)]))  # Correction here
    balance=defaultdict(lambda: defaultdict(int))

    if records ==[]:
        return balance
    # If the retrieved records are none, return a defaultdict of defaultdict of integer
    if not records:
        return balance
    
    
    for record in records:
        if record['data']['method']=='mint':
            balance[record['data']['ticker']][record['from_address']]=record['data']['amount']
        if record['data']['method']=='transfer':
            balance[record['data']['ticker']][record['from_address']]-=record['data']['amount']
            balance[record['data']['ticker']][record['data']['to']]+=record['data']['amount']

    return balance


def write_to_mongodb(block_no,KZG_Hash,BS20_json,from_address):
    # Assuming 'credentials' contains MongoDB connection info
    client = MongoClient(credential['mongo_uri'])  # Use your MongoDB URI
    db = client.BlobScriptions  # Access the database using the client
    collection = db.BlobScriptions  # Access the collection from the database

    token_balance = track_balance()

    if BS20_json['method']=='mint':
        if collection.find_one({"KZG_Hash": KZG_Hash}):
            print("KZG commit is a duplicate, not inserting.")
        if BS20_json['ticker'] in token_balance.keys():
            print("Token With Same Ticker Already Exists")
            return None
        else:
            # Prepare the document to insert
            document = {
                "block_no": block_no,
                "KZG_Hash": KZG_Hash,
                "from_address": from_address,
                "data": BS20_json
            }
            collection.insert_one(document)
            print("Document inserted successfully.")
        #collection.insert_one(BS20_json)  # Insert the JSON data into MongoDB

    if  BS20_json['method']=='transfer':
        #check balance
        ticker = BS20_json.get('ticker')
        from_address_balance = token_balance.get(ticker, {}).get(from_address, 0)
        if from_address_balance < BS20_json.get('amount', 0):
            print("Not Enough Balance")
            return None
        else:
            if collection.find_one({"KZG_Hash": KZG_Hash}):
                print("KZG commit is a duplicate, not inserting.")
            else:
                # Prepare the document to insert
                document = {
                    "block_no": block_no,
                    "KZG_Hash": KZG_Hash,
                    "from_address": from_address,
                    "data": BS20_json
                }
                collection.insert_one(document)
                print("Document inserted successfully.")

start_index_dummy()