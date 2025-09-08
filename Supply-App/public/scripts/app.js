async function submitBatch() {
    const form = document.getElementById("batchForm");

    const data = {
      BatchID: form.BatchID.value,
      Type: form.Type.value,
      Quantity: parseInt(form.Quantity.value),
      ManufactureDate: form.ManufactureDate.value,
      ExpiryDate: form.ExpiryDate.value,
      Status: form.Status.value,
      Composition: form.Composition.value,
      Inspection: form.Inspection.value,
      Serials: form.Serials.value
    };

    try {
      const response = await fetch("/api/batch", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
      });

      if (!response.ok) {
        const errText = await response.text();
        alert("Error:\n" + errText);
        return;
      }

      const result = await response.json();

      // Build formatted alert message
      let msg = "✅ Batch Created Successfully!\n\n";
      for (let key in result) {
        msg += `${key}: ${result[key]}\n`;
      }
      alert(msg);

      form.reset();
    } catch (err) {
      alert("Fetch error: " + err);
    }
  }

  // Query Batch
  const readData = async (event) => {
    event.preventDefault();
    const batchIDInput = document.getElementById("batchIdInput").value;

    if (batchIDInput.trim().length === 0) {
      alert("⚠️ Please enter a valid ID.");
    } else {
      try {
        const response = await fetch(`/api/batch/${batchIDInput}`);
        let responseData = await response.json();

        if (!response.ok) {
          alert("Error:\n" + JSON.stringify(responseData));
          return;
        }

        // Build formatted alert message  
        let msg = "📦 Batch Details:\n\n";
        for (let key in responseData) {
          msg += `${key}: ${responseData[key]}\n`;
        }
        alert(msg);

      } catch (err) {
        alert("❌ Error fetching batch.");
        console.log(err);
      }
    }
  };