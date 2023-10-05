$(document).ready(function() {
    $.get("/order/ids", function(data) {
        if (data.status === "OK" && data.ids.length > 0) {
            var orderDropdownMenu = $("#orderDropdownMenu");
            data.ids.forEach(function(orderId) {
                var menuItem = $('<li><a class="dropdown-item" href="#">' + orderId + '</a></li>');
                menuItem.click(function() {
                    loadOrder(orderId);
                });
                orderDropdownMenu.append(menuItem);
            });
        } else {
            console.log("Не удалось получить id заказов или список пуст.");
        }
    });
});

function loadOrder(orderUid) {
    $.get("/order/get?order_uid=" + orderUid, function(data) {
        $("#orderIdHeader").text("Заказ № " + orderUid);
        var dateCreated = new Date(data.order.date_created);
        var formattedDate = formatDate(dateCreated);
        $("#orderDate").text(formattedDate);
        $("#deliveryService").text(data.order.delivery_service);
        $("#recipientName").text(data.order.delivery.name);
        $("#recipientPhone").text(data.order.delivery.phone);
        $("#recipientAddress").text(data.order.delivery.address + ", " + data.order.delivery.city + ", " + data.order.delivery.region + ", " + data.order.delivery.zip);
        $("#recipientEmail").text(data.order.delivery.email);
        $("#paymentBank").text(data.order.payment.bank);
        $("#paymentProvider").text(data.order.payment.provider);
        $("#goodsTotal").text(data.order.payment.goods_total + " " + data.order.payment.currency);
        $("#deliveryCost").text(data.order.payment.delivery_cost + " " + data.order.payment.currency);
        $("#customFee").text(data.order.payment.custom_fee + " " + data.order.payment.currency);
        $("#totalAmount").text(data.order.payment.amount + " " + data.order.payment.currency);

        $("#orderedItems").empty();
        data.order.items.forEach(function(item) {
            var row = '<tr>' +
                '<td>' + item.name + ' </td>' +
                '<td>' + item.price + ' ' + data.order.payment.currency + '</td>' +
                '<td>' + item.sale + '%</td>' +
                '<td>' + item.size + '</td>' +
                '<td>' + item.brand + '</td>' +
                '<td>' + item.total_price + ' ' + data.order.payment.currency + '</td>' +
                '</tr>';
            $("#orderedItems").append(row);
        });
    });
}

function formatDate(date) {
    var options = { year: 'numeric', month: 'long', day: 'numeric' };
    var formattedDate = date.toLocaleDateString('ru-RU', options);
    return formattedDate;
}