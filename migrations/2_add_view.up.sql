CREATE OR REPLACE VIEW public.orders_view
AS SELECT ord.order_uid,
    ord.track_number,
    ord.entry,
    ord.locale,
    ord.internal_signature,
    ord.customer_id,
    ord.delivery_service,
    ord.sm_id,
    ord.date_created,
    del.name AS delivery_name,
    del.phone AS delivery_phone,
    del.zip AS delivery_zip,
    del.city AS delivery_city,
    del.address AS delivery_address,
    del.region AS delivery_region,
    del.email AS delivery_email,
    pay.transaction AS payment_transaction,
    pay.request_id AS payment_request_id,
    pay.currency AS payment_currency,
    pay.provider AS payment_provider,
    pay.amount AS payment_amount,
    pay.payment_dt AS payment_payment_dt,
    pay.bank AS payment_bank,
    pay.delivery_cost AS payment_delivery_cost,
    pay.goods_total AS payment_goods_total,
    pay.custom_fee AS payment_custom_fee,
    array_agg(i.*) AS items
   FROM orders ord
     JOIN deliveries del ON ord.delivery_id = del.delivery_id
     JOIN payments pay ON ord.order_uid = pay.transaction
     JOIN items i ON ord.track_number::text = i.track_number::text
  GROUP BY ord.order_uid, ord.track_number, ord.entry, ord.locale, ord.internal_signature, ord.customer_id, ord.delivery_service, ord.sm_id, ord.date_created, del.name, del.phone, del.zip, del.city, del.address, del.region, del.email, pay.transaction, pay.request_id, pay.currency, pay.provider, pay.amount, pay.payment_dt, pay.bank, pay.delivery_cost, pay.goods_total, pay.custom_fee;-