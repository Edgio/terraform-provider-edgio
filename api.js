
async function getEdgioToken() {
    const clientId = ''; 
    const clientSecret = '';
    const url = 'https://id.edgio.app/connect/token';
    const data = {
      client_id: clientId,
      client_secret: clientSecret,
      grant_type: 'client_credentials',
      scope: 'app.accounts'
    };
  
    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      },
      body: new URLSearchParams(data) // Encodes data as URL-encoded form data
    };
  
    try {
      const response = await fetch(url, options);
      if (!response.ok) {
        throw new Error(`Error fetching token: ${response.status}`);
      }
      const responseData = await response.json();

      return responseData.access_token;
    } catch (error) {
      console.error('Error:', error);
    }
  }
  
  async function getEdgioProperties(bearerToken, page, pageSize, organizationId) {
    const url = new URL('https://edgioapis.com/accounts/v0.1/properties');
    url.searchParams.set('page', page); // Add page parameter
    url.searchParams.set('page_size', pageSize); // Add page size parameter
    url.searchParams.set('organization_id', organizationId); // Add organization ID
  
    const options = {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${bearerToken}`
      }
    };
  
    try {
      const response = await fetch(url, options);
      
      const responseData = await response.json();
      console.log('Edgio Properties:', responseData);
    } catch (error) {
      console.error('Error:', error);
    }
  }

async function DisplayProperties() {
    const bearerToken = await getEdgioToken();
    const page = 1;
    const pageSize = 10;
    const organizationId = ''; // Replace with your actual organization ID
    const p = await getEdgioProperties(bearerToken, page, pageSize, organizationId);

    console.log(p);
}

DisplayProperties();