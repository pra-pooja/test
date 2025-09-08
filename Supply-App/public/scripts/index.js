const addData = async (event) => {
    event.preventDefault();
  
    const BatchID = document.getElementById("BatchID").value;
    const Type = document.getElementById("Type").value;
    const Quantity = document.getElementById("Quantity").value;
    const ManufactureDate = document.getElementById("ManufactureDate").value;
    const ExpiryDate = document.getElementById("ExpiryDate").value;
    const Status = document.getElementById("Status").value;
    const Composition = document.getElementById("Composition").value;
    const Inspection = document.getElementById("Inspection").value;
    const Serials = document.getElementById("Serials").value;

    const batchData = {
      BatchID: BatchID,
      Type: Type,
      Quantity: Quantity,
      ManufactureDate: ManufactureDate,
      ExpiryDate: ExpiryDate,
      Status: Status,
      Composition: Composition,
      Inspection: Inspection,
      Serials: Serials
    };
  
    if (
      BatchID.length == 0 ||
      Type.length == 0 ||
      Quantity.length == 0 ||
      ManufactureDate.length == 0 ||
      ExpiryDate.length == 0 ||
      Status.length == 0  ||
      Composition.length == 0 ||
      Inspection.length == 0 |
      Serials.length == 0
    ) {
      alert("Please enter the data properly.");
    } else {
      try {
        const response = await fetch("/api/batch", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(batchData),
        });
        console.log("RESPONSE: ", response);
        const data = await response.json();
        console.log("DATA: ", data);
        return alert("Batch Created");
      } catch (err) {
        alert("Error");
        console.log(err);
      }
    }
  };
  
  const readData = async (event) => {
    event.preventDefault();
    const batchIDInput = document.getElementById("batchIdInput").value;
    console.log(batchIDInput)
  
    if (batchIDInput.length == 0) {
      alert("Please enter a valid ID.");
    } else {
      try {
        const response = await fetch(`/api/batch/${batchIDInput}`);
        let responseData = await response.json();
        console.log("response data", responseData);
        alert(JSON.stringify(responseData));
      } catch (err) {
        alert("Error");
        console.log(err);
      }
    }
  };
  