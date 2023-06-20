// Функция для получения списка заказов
function getOrders() {
    fetch('http://localhost:8080/get_orders', {
      credentials: 'include',
    })
      .then(response => response.json())
      .then(function(responseData) {
        console.log(responseData)
        var data = responseData;
        const orderList = document.getElementById('order-list');
  
        orderList.innerHTML = '';
  
        data.order.forEach(order => {
          const listItem = document.createElement('li');
          listItem.textContent = `Product: ${order.ProductName}, Quantity: ${order.ProductQuantity}`;
          orderList.appendChild(listItem);
        });
      })
      .catch(error => {
        alert("You should sign in");
        window.location.href = '../index.html';
        console.error('Ошибка при получении списка заказов:', error);
      });
  }
  
  // Обработчик клика по кнопке "Вернуться"
  function goBack() {
    window.location.href = '../products/products.html'; // Замените на свой URL страницы с продуктами
  }
  
  // Получаем список заказов при загрузке страницы
  document.addEventListener('DOMContentLoaded', getOrders);
  
  // Назначаем обработчик клика по кнопке "Вернуться"
  const backButton = document.getElementById('back-button');
  backButton.addEventListener('click', goBack);
  
  // Получение ссылки на кнопку выхода из аккаунта
const logoutBtn = document.getElementById('logoutBtn');

// Обработчик события клика на кнопку выхода из аккаунта
logoutBtn.addEventListener('click', logout);

// Функция для выхода из аккаунта
function logout() {
  fetch('http://localhost:8080/sign_out',
  {
      credentials: 'include',
  })
  .then(function(response) {
    if (response.ok) {
      setCookie("session_token", "", {"max-age":0})
      window.location.href = '../index.html';
    }
    throw new Error('Ошибка при выполнении запроса.');
  })
}