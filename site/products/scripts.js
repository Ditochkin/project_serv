function loadProducts() {
    fetch('http://localhost:8080/get_products',
    {
        credentials: 'include',
    })
      .then(function(response) {
        if (response.ok) {
          return response.json();
        }
        alert("You should sign in");
        window.location.href = '../index.html';

        throw new Error('Ошибка при выполнении запроса.');
      })
      .then(function(responseData) {
        console.log(responseData)
        var productsContainer = document.getElementById('products');
  
        responseData.products.forEach(function(product) {
            var productElement = document.createElement('div');
            productElement.classList.add('product');
          
            productElement.innerHTML = `
                <h3>${product.ProductName}</h3>
                <p class="product-description">${product.ProductDescription}</p>
                <p class="product-price">Price: ${product.ProductPrice}</p>
                <p class="product-quantity">Quantity: ${product.ProductQuantity}</p>
                <div class="product-quantity">
                  <label class="quantity-label">Quantity to buy:</label>
                  <input class="quantity-input" type="number" min="0" value="0">
                  <button class="buy-button">Buy</button>
                </div>
            `;
          
            var buyButton = productElement.querySelector('.buy-button');
            var quantityInput = productElement.querySelector('.quantity-input');
            var productQuantityElement = productElement.querySelector('.product-quantity');
          
            buyButton.addEventListener('click', function() {
                var quantity = parseInt(quantityInput.value);
                var currentQuantity = parseInt(productQuantityElement.textContent.split(': ')[1]);
        
                if (quantity > currentQuantity) {
                    alert('The number of products in the order exceeds the available quantity!');
                    return
                }
              
                if (quantity > 0) {
                    console.log(`Купить ${quantity} шт. продукта "${product.ProductName}"`);
                
                    var productData = {
                        "productname": product.ProductName,
                        "productquantity": quantity
                    };
                
                    fetch('http://localhost:8080/add_order', {
                        credentials: 'include',
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify(productData)
                    })
                    .then(function(response) {
                        if (response.ok) {
                            console.log('Заказ успешно добавлен');
                        } else {
                            console.log('Ошибка при добавлении заказа');
                        }
                    })
                    .catch(function(error) {
                        console.log('Ошибка при выполнении запроса', error);
                    });
        
                    var updatedQuantity = currentQuantity - quantity;
                    productQuantityElement.textContent = `Quantity: ${updatedQuantity}`;

                    quantityInput.value = 0;

                    showPopup(`You bought ${quantity} pieces of the "${product.ProductName}"`, 2000);
                }
            });
          
            productsContainer.appendChild(productElement);

        });
      })
      .catch(function(error) {
        console.error(error);
      });
  }
  
  function buyProduct(productName, productQuantity) {
    var cartCountElement = document.getElementById('cart-count');
    var currentCount = parseInt(cartCountElement.innerText);
    cartCountElement.innerText = currentCount + productQuantity;
  }

  document.addEventListener('DOMContentLoaded', function() {
    loadProducts();
  });

  const cartIcon = document.getElementById('cart-icon');

  cartIcon.addEventListener('click', () => {
    window.location.href = '../bascket/bascket.html';
  });

  // Получение ссылки на кнопку выхода из аккаунта
const logoutBtn = document.getElementById('logoutBtn');

// Обработчик события клика на кнопку выхода из аккаунта
logoutBtn.addEventListener('click', logout);


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