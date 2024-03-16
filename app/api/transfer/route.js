export async function POST(req) {
  console.log("handle POST transfer")
  try {
    // Extract mintData from the request body
    const { transferData } = req.body;

    // Process mintData (e.g., mint tokens)
    // Your minting logic goes here...

    // If minting is successful, send a success response
    // return new Response(
    //   {message: 'server:Token minted successfully',},
    //   {status: 200}
    // )

    return Response.json(
      {message:"server:Token transferred successfully"},
    {status: 200,
      
    })
  } catch (error) {
    console.log(error)
    return Response.json(
      {message:  error.message|"Internal server error",},
      {status: 200,
      }
    )

  }
  }