package orchestrator

import "net/http"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Calculator API Interface</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      padding: 2rem;
      background: #f9f9f9;
    }
    .container {
      max-width: 600px;
      margin: auto;
      background: #fff;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 0 10px rgba(0,0,0,0.1);
    }
    input, button {
      margin: 8px 0;
      padding: 8px;
      font-size: 1rem;
    }
    .result {
      margin: 10px 0;
      padding: 10px;
      border: 1px solid #ccc;
      background: #eee;
      border-radius: 4px;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Calculator API Interface</h1>
    
    
    <form id="calculateForm">
      <label for="expression">Enter an expression:</label><br>
      <input type="text" id="expression" name="expression" placeholder="2 * 2" required><br>
      <button type="submit">Send</button>
    </form>
    <div id="taskResult" class="result" style="display: none;"></div>
    
    <hr>
    
    
    <div id="finalResult" class="result" style="display: none;"></div>
  </div>
  
  <script>
    async function pollResult(taskId) {
      const resultDiv = document.getElementById('finalResult');
      resultDiv.style.display = 'block';
      resultDiv.innerText = "Calculated...";
      try {
        const response = await fetch('http://localhost:8080/api/v1/expressions/' + encodeURIComponent(taskId));
        if (!response.ok) {
          resultDiv.innerText = 'Error getting result: ' + await response.text();
          return;
        }
        const data = await response.json();
        const expr = data.expression;
        if (expr.status === "done") {
          resultDiv.innerText = 'Result: ' + expr.result;
        } else if (expr.status === "failed") {
          resultDiv.innerText = 'Calculation error: ' + expr.ErrorMessage;
        } else {
          setTimeout(() => pollResult(taskId), 1000);
        }
      } catch (error) {
        resultDiv.innerText = 'Request error: ' + error;
      }
    }

    document.getElementById('calculateForm').addEventListener('submit', async function(e) {
      e.preventDefault();
      const finalResultDiv = document.getElementById('finalResult');
      finalResultDiv.style.display = 'block';
      finalResultDiv.innerText = "Calculated...";
      
      try {
        const response = await fetch('http://localhost:8080/api/v1/calculate', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ expression: document.getElementById('expression').value })
        });
        const taskResultDiv = document.getElementById('taskResult');
        if (!response.ok) {
          const errorData = await response.json();
          taskResultDiv.style.display = 'block';
          taskResultDiv.innerText = 'Error: ' + errorData.error_message;
          return;
        }
        const resJson = await response.json();
        taskResultDiv.style.display = 'block';
        taskResultDiv.innerText = 'Expression has been set, ID: ' + resJson.id;
        pollResult(resJson.id);
      } catch (error) {
        const taskResultDiv = document.getElementById('taskResult');
        taskResultDiv.style.display = 'block';
        taskResultDiv.innerText = 'Request error: ' + error;
      }
    });
  </script>
</body>
</html>`
	w.Write([]byte(html))
}
