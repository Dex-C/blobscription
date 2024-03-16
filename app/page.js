"use client"
import Image from "next/image";
import { useState } from 'react';
import styles from "./page.module.css";
import { ethers } from "ethers";

export default function Home() {
  const [mintFormData, setMintFormData] = useState({
    tickit: '',
    to: '',
    amount: 0,
    img: ''
  });

  const [transferFormData, setTransferFormData] = useState({
    tickit: '',
    to: '',
    amount: 0,
    img: ''
  });


  // // for popout error warnings
  const [error, setError] = useState('');

  // handle form changes
  const handleMintChange = (e) => {
    const { name, value } = e.target;
    setMintFormData({ ...mintFormData, [name]: value });
  };
  const handleTransferChange = (e) => {
    const { name, value } = e.target;
    setTransferFormData({ ...transferFormData, [name]: value });
  };

  const handleMint = async (e) => {
    e.preventDefault();

    // Clear any previous error
    setError('');

    try {
      // Optionally, show a success message to the user
      const message = `Form submitted successfully!\n\nTickit: ${mintFormData.tickit}\nTo: ${mintFormData.to}\nAmount: ${mintFormData.amount}\nImage: ${mintFormData.img}`;
      alert(message)

      // Perform form submission here, for example, sending data to a server
      let mintData = mintFormData
      mintData.method = "mint"
      const response = await fetch('/api/mint', {
        method: 'POST',
        body: JSON.stringify(mintData),
        headers: {
          'Content-Type': 'application/json'
        }
      });
      console.log(response)
      console.log(response.body)
      const data = await response.json();
      console.log(data)

      // If there's an error returned from the server, set the error state
      if (!response.ok) {
        console.log("not ok")
        throw new Error(data.message);
      } else {
        console.log("ok")
        alert("token minted successfully")
      }

      // Reset form fields after successful submission
      setMintFormData({
        tickit: '',
        to: '',
        amount: 0,
        img: ''
      });


    } catch (error) {
      // Handle and display any errors to the user
      console.log(error)
      setError(error.message || 'An error occurred while submitting the form.');
    }
  };

  const handleTransfer = async (e) => {
    e.preventDefault();

    // Clear any previous error
    setError('');

    try {
      // Optionally, show a success message to the user
      const message = `Form submitted successfully!\n\nTickit: ${transferData.tickit}\nTo: ${transferData.to}\nAmount: ${transferData.amount}\nImage: ${transferData.img}`;
      alert(message)

      // Perform form submission here, for example, sending data to a server
      let transferData = transferFormData
      transferData.method = "transfer"
      const response = await fetch('/api/transfer', {
        method: 'POST',
        body: JSON.stringify(transferData),
        headers: {
          'Content-Type': 'application/json'
        }
      });
      const data = await response.json();
      console.log(data)

      // If there's an error returned from the server, set the error state
      if (!response.ok) {
        console.log("not ok")
        throw new Error(data.message);
      } else {
        console.log("ok")
        alert("token transferred successfully")
      }

      // Reset form fields after successful submission
      setTransferFormData({
        tickit: '',
        to: '',
        amount: 0,
        img: ''
      });


    } catch (error) {
      // Handle and display any errors to the user
      console.log(error)
      setError(error.message || 'An error occurred while submitting the form.');
    }
  };
  return (
    <main className={styles.main}>
      <h3>wallet address {process.env.NEXT_PUBLIC_ADDRESS}</h3>
      <></>
      {/* {mint} */}
      <h3>Mint Tokens</h3>
      {error && <div style={{ color: 'red' }}>{error}</div>}
      <form onSubmit={handleMint}>
        <div>
          <label>Tickit:</label>
          <input type="text" name="tickit" value={mintFormData.tickit} onChange={handleMintChange} />
        </div>
        <div>
          <label>To:</label>
          <input type="text" name="to" value={mintFormData.to} onChange={handleMintChange} />
        </div>
        <div>
          <label>Amount:</label>
          <input type="number" name="amount" value={mintFormData.amount} onChange={handleMintChange} />
        </div>
        <div>
          <label>Image URL:</label>
          <input type="text" name="img" value={mintFormData.img} onChange={handleMintChange} />
        </div>
        <button type="submit">Mint</button>
      </form>

      {/* {transfer} */}

      <h3>Transfer Tokens</h3>
      {error && <div style={{ color: 'red' }}>{error}</div>}
      <form onSubmit={handleTransfer}>
        <div>
          <label>Tickit:</label>
          <input type="text" name="tickit" value={transferFormData.tickit} onChange={handleTransferChange} />
        </div>
        <div>
          <label>To:</label>
          <input type="text" name="to" value={transferFormData.to} onChange={handleTransferChange} />
        </div>
        <div>
          <label>Amount:</label>
          <input type="number" name="amount" value={transferFormData.amount} onChange={handleTransferChange} />
        </div>
        <div>
          <label>Image URL:</label>
          <input type="text" name="img" value={transferFormData.img} onChange={handleTransferChange} />
        </div>
        <button type="submit">Transfer</button>
      </form>

      <h2>Dashboard</h2>
    </main>


  );
}
