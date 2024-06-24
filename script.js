document.getElementById('houseForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const formData = {
        size: parseFloat(document.getElementById('size').value),
        bedrooms: parseFloat(document.getElementById('bedrooms').value),
        age: parseFloat(document.getElementById('age').value),
        location: document.getElementById('location').value
    };

    try {
        const response = await fetch('http://localhost:8000/predict', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ house: formData })
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const data = await response.json();
        document.getElementById('prediction').innerText = `Precio predicho: $${(data.Price.toFixed(2))}`;
        //console.log(data.price)
    
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('prediction').innerText = 'Error predicting price';
    }
});
