let isRegistered = 0;

function showPopup(message, duration) {
  var popup = document.createElement("div");
  popup.className = "popup";
  popup.textContent = message;
  document.body.appendChild(popup);

  setTimeout(function() {
    document.body.removeChild(popup);
  }, duration);
}

function login() {
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
  
    var data = {
      username: username,
      password: password
    };

    console.log(data)
  
    var url = new URL("http://localhost:8080/sign_in");
  
    fetch(url, {
      credentials: 'include',
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Content-Length": JSON.stringify(data).length.toString()
      },
      body: JSON.stringify(data)
    })
      .then(function(response) {
        console.log(response)
        if (response.ok) {
            console.log("OK")
            window.location.href = 'products/products.html';
          return response.json();
        }
        if (response.status == 500)
        {
          showPopup("Login or password is incorrect!", 3000);
        }
        throw new Error("Ошибка при выполнении запроса.");
      })
      .then(function(responseData) {
        // Обработка ответа от сервера
        console.log(responseData);
        setCookie("session_token", responseData.tokens.token, {"max-age":3600})
      })
      .catch(function(error) {
        console.error(error);
      });
  }

  function signup() {
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    var email = document.getElementById("email").value;

    console.log(username, password, email)

    if (username == "")
    {
      showPopup("Login is empty!", 3000);
      return;
    }
    if (password == "")
    {
      showPopup("Password is empty!", 3000);
      return;
    }
    if (email == "")
    {
      showPopup("E-mail is empty!", 3000);
      return;
    }
  
    var data = {
      username: username,
      password: password,
      email: email,
    };

    console.log(data)
  
    var url = new URL("http://localhost:8080/create_user");
  
    fetch(url, {
      credentials: 'include',
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Content-Length": JSON.stringify(data).length.toString()
      },
      body: JSON.stringify(data)
    })
      .then(function(response) {
        console.log(response)
        if (response.ok) {
            console.log("OK")
            window.location.href = 'index.html';
          return;
        }
        if (response.status == 500)
        {
          showPopup("This login or email already exists!", 3000);
        }
        throw new Error("Ошибка при выполнении запроса.");
      })
      .catch(function(error) {
        console.error(error);
      });
  }


function toggleFormSignIn() {
    var loginForm = document.getElementById("login-form");
    var registerForm = document.getElementById("register-form");
  
    loginForm.style.display = loginForm.style.display === "none" ? "block" : "none";
    registerForm.style.display = registerForm.style.display === "none" ? "block" : "none";
  }

function toggleFormSignUp() {
    var loginForm = document.getElementById("login-form");
    var registerForm = document.getElementById("register-form");
  
    loginForm.style.display = loginForm.style.display === "none" ? "block" : "none";
    registerForm.style.display = registerForm.style.display === "none" ? "block" : "none";
  }
  
document.addEventListener("DOMContentLoaded", function() {
  var loginForm = document.getElementById("login-form");
  loginForm.style.display = "block";
});