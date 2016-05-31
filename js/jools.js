"use strict;";

Cart = function() {
    var MAX_PRODUCTS_IN_CART = 50;
    var baseKey = "myCart";
    var syncId = 2;
    this.getBaseKey = function() {
        return baseKey;
    };

    this.getSyncUrl = function() {
        return "/cart";
    };
    this.isTemp = function() {
        return false;
    };
    this.MergeForSync = function(localData, tempData) {
      localData = JSON.parse(localData);
      tempData = JSON.parse(tempData);
      var mergedData = this.mergeCarts(localData, tempData);
      return JSON.stringify(mergedData);
    };
    
    this.ByDecreasingPrice = function(a,b) {
      return a.Price < b.Price;
    };
    this.mergeCarts = function(cart1, cart2) {
      //A merge of cart1 and cart2 does not add quantities of same prodIndex. 
      //It overrides the quantities mentioned in cart1 by those of cart2.
      var mergedCart = this.getEmptyCart();
      if (cart1.NumProducts == 0) {
        return cart2;
      }
      if (cart2.NumProducts == 0) {
        return cart1;
      }

      var index1 = 0, index2 = 0, indexM = 0;
      for(;index1 < cart1.Products.length && index2 < cart2.Products.length; indexM++) {
        if(cart1.Products[index1].Index == cart2.Products[index2].Index) {
          //Comparing indices for equality as prices are floats
          mergedCart.Products[indexM] = cart2.Products[index2];
          mergedCart.Indices[mergedCart.Products[indexM].Index] = indexM;
          index1++;
          index2++;
        } else if(cart1.Products[index1].Price > cart2.Products[index2].Price) {
          mergedCart.Products[indexM] = cart1.Products[index1];
          mergedCart.Indices[mergedCart.Products[indexM].Index] = indexM;
          index1++;
        } else {
          mergedCart.Products[indexM] = cart2.Products[index2];
          mergedCart.Indices[mergedCart.Products[indexM].Index] = indexM;
          index2++;
        }
      } 
      if (index1 < cart1.Products.length) {
        for(;index1 < cart1.Products.length;index1++, indexM++) {
          mergedCart.Products[indexM] = cart1.Products[index1];
          mergedCart.Indices[mergedCart.Products[indexM].Index] = indexM;
        }
      }
      if (index2 < cart2.Products.length) {
        for(;index2 < cart2.Products.length;index2++, indexM++) {
          mergedCart.Products[indexM] = cart2.Products[index2];
          mergedCart.Indices[mergedCart.Products[indexM].Index] = indexM;
        }
      }
      mergedCart.NumProducts = indexM;
      return mergedCart;
    };
    this.getEmptyCart = function() {
        return {NumProducts: 0, Indices:{}, Products: []};
    };
    this.getProducts = function() {
        var storageKey = baseKey;
        var myCart = localStorage.getItem(storageKey);
        if (myCart === null || myCart == "") {
          myCart = this.getEmptyCart();
        } else {
          myCart = JSON.parse(myCart);
        }

        storageKey = utilObj.getLocalKey(baseKey);
        var myLocalCart = localStorage.getItem(storageKey);
        if (myLocalCart === null || myCart == "") {
          myLocalCart = this.getEmptyCart();
        } else {
          myLocalCart = JSON.parse(myLocalCart);
        }

        storageKey = utilObj.getTempKey(baseKey);
        var myTempCart = localStorage.getItem(storageKey);
        if (myTempCart === null || myTempCart == "") {
          myTempCart = this.getEmptyCart();
        } else {
          myTempCart = JSON.parse(myTempCart);
        }

        var allProducts = this.mergeCarts(myCart, myLocalCart);
        if (myTempCart.NumProducts != 0) {
          allProducts = this.mergeCarts(allProducts, myTempCart);
        }
        return allProducts;
    };
    this.getRemovedProducts = function() {
        var storageKey = utilObj.getRemoveKey(baseKey);
        var removed = localStorage.getItem(storageKey);
        if (removed === null || removed == "") {
          removed = this.getEmptyCart();
        } else {
          removed = JSON.parse(removed);
        }
        return removed;
    };
    this.isProductInCart = function(prodIndex) {
        var allProducts = this.getProducts();
        if (allProducts.Indices.hasOwnProperty(prodIndex)) {
            return true;
        }
        return false;
    };
    this.addToCart = function(productInfo, showMsg) {
        if (productInfo == undefined) {
          console.log("Product to be added to cart is undefined");
          return false;
        }
        var storageKey = utilObj.getLocalKey(baseKey);
        if (syncObj.isSyncInProgress(syncId)) {
            storageKey = utilObj.getTempKey(baseKey);
        }
        var curCart = localStorage.getItem(storageKey);
        if (curCart === null) {
          curCart = this.getEmptyCart();
        } else {
          curCart = JSON.parse(curCart);
        }

        if (curCart.Indices.hasOwnProperty(productInfo.Index)) {
            //Duplicate item entry
            curCart.Products[ curCart.Indices[productInfo.Index] ].Qty++;
        } else {
            //New item being added 
            if(curCart.NumProducts >= MAX_PRODUCTS_IN_CART) {
              utilObj.showMessage("There are too many products in your shopping bag. Unable to add more.", false);
              return false;
            }
            curCart.Products[curCart.NumProducts] = productInfo;
            curCart.NumProducts++;
            curCart.Products.sort( cartObj.ByDecreasingPrice );
            //Rebuild cart Index
            var rebuiltIndex = {};
            for(var i = 0; i < curCart.NumProducts; i++) {
              rebuiltIndex[ curCart.Products[i].Index ] = i;
            }
            curCart.Indices = rebuiltIndex;
        }
    
        this.setProducts(storageKey, curCart);
        this.updateItemCount();
        if(showMsg == true) {
            utilObj.showMessage("Successfully added a " + productInfo.Desc + " to your shopping bag.", true);
        }
        
        storageKey = utilObj.getRemoveKey(baseKey);
        this.checkAndRemove(storageKey, productInfo.Index);
        return true;
    };
    this.setProducts = function(storageKey, Products) {
        localStorage.setItem(storageKey, JSON.stringify(Products));
    };
    this.removeFromCart = function(prodIndex) {
        if(prodIndex == undefined) {
            console.log("Index of product to be removed is not defined");
            return false;
        }
        if(!this.isProductInCart(prodIndex)) {
            return false;
        }
        
        var removedItems = this.getRemovedProducts();
        if (removedItems.Indices.hasOwnProperty(prodIndex)) {
            console.log("Attempted to remove again", prodIndex);
            return false;
        }
        var storageKey = utilObj.getRemoveKey(baseKey);
        removedItems.Indices[prodIndex] = removedItems.NumProducts;
        removedItems.NumProducts++;
        this.setProducts(storageKey, removedItems);

        storageKey = utilObj.getTempKey(baseKey);
        this.checkAndRemove(storageKey, prodIndex);

        storageKey = utilObj.getLocalKey(baseKey);
        this.checkAndRemove(storageKey, prodIndex);

        storageKey = baseKey;
        this.checkAndRemove(storageKey, prodIndex);

        this.updateItemCount();
        this.showItemList(true, null);
        return true;
    };
    this.checkAndRemove = function(storageKey, prodIndex) {
        //What is required here is to remove the item from the list and not worry about 
        //reducing quantity etc
        var curItems = localStorage.getItem(storageKey);
        if (curItems === null) {
          return false;
        }
        curItems = JSON.parse(curItems);
        if (!curItems.Indices.hasOwnProperty(prodIndex)) {
          return false;
        }
        var updatedItems = this.getEmptyCart();
        var Products = curItems.Products;
        
        for(var i=0; i < Products.length; i++) {
            if(Products[i].Index != prodIndex) {
                updatedItems.Products[updatedItems.NumProducts] = Products[i];
                updatedItems.Indices[Products[i].Index] = updatedItems.NumProducts;
                updatedItems.NumProducts++;
            }
        }
        if(Products.length == 0) {
            //To handle the case of removedProducts, which doesn't have product list
            delete curItems.Indices[prodIndex];
            curItems.NumProducts--;
            this.setProducts(storageKey, curItems);
        } else {
            this.setProducts(storageKey, updatedItems);
        }
        return true;
    };
    
    this.updateQuantity = function(prodIndex, qty) {
        qty = parseInt(qty);
        if (qty > 10) {
            console.log("Update qty called with qty > 10");
            return false;
        }
        if(!this.isProductInCart(prodIndex)) {
            console.log("Update qty called for product not in cart");
            return false;
        }
        if(qty <= 0) {
            return this.removeFromCart(prodIndex);
        }
        var storageKey = baseKey;
        var myCart = localStorage.getItem(storageKey);
        if (myCart === null) {
          myCart = this.getEmptyCart();
        } else {
          myCart = JSON.parse(myCart);
        }
        var productData = null;
        if (myCart.Indices.hasOwnProperty(prodIndex)) {
          productData = myCart.Products[ myCart.Indices[prodIndex] ];
          this.checkAndRemove(storageKey, prodIndex);
        }
        storageKey = utilObj.getLocalKey(baseKey);
        if (syncObj.isSyncInProgress(syncId)) {
            storageKey = utilObj.getTempKey(baseKey);
        }
        var localCart = localStorage.getItem(storageKey);
        if (localCart === null) {
          localCart = this.getEmptyCart();
        } else {
          localCart = JSON.parse(localCart);
        }
        if (localCart.Indices.hasOwnProperty(prodIndex)) {
          localCart.Products[ localCart.Indices[prodIndex] ].Qty = qty;
          this.setProducts(storageKey, localCart);
        } else {
          productData.Qty = qty;
          this.addToCart(productData, false);
        }
        this.showItemList(true, null);
        return true;
    };

    this.getNumProducts = function() {
        var cartData = this.getProducts();
        return cartData.NumProducts;
    };
    this.updateItemCount = function() {
        var numProducts = cartObj.getNumProducts();
        $('.cart_item_count').text(numProducts);
    };
    
    this.showItemList = function(isEditable, productList) {
        var allProducts;
        if(productList === null) {
          allProducts = this.getProducts();
        } else {
          allProducts = productList;
        }
        var Products = allProducts.Products;
        var markup = '';
        var total = 0;
        if(Products == null || Products.length == 0) {
            markup = "Your shopping bag is currently empty";
        } else {  
            for(var i = 0; i < Products.length; i++) {
                var subTotal = Products[i].Price * Products[i].Qty;
                markup += '<div class="bag_item_wrap">'
                        +   '<img src="' + configObj.DESIGNS_URL + '/' + Products[i].Image + '" alt="" width="120" class="left" style="margin-right: 20px;">'
                        +   '<div class="bag_item_desc left">'
                        +       '<div class="bag_item_name">'
                        +           Products[i].Name
                        +       '</div>'
                        +       '<div>'
                        +           Products[i].Desc
                        +       '</div>';
                if(isEditable) {
                  markup +=    '<div class="bag_item_actions">'
                        +           '<span class="remove_item like_link" data-index="' + Products[i].Index + '">remove</span>'
                        +       '</div>';
                }
                markup  +=  '</div>'
                        +   '<div class="left center" style="width: 160px;margin-top: 15px;">'
                        +       utilObj.formatMoney(Products[i].Price, true)
                        +   '</div>';
                var qtyMarkup = '<select class="cart_qty" data-index="' + Products[i].Index + '">';
                var qtyMatched = false;
                for (var j = 1; j <= 10; j++) {
                    var selected = '';
                    qtyMatched = false;
                    if(Products[i].Qty == j) {
                        selected = 'selected="selected"';
                        qtyMatched = true;
                    }
                    if(isEditable) {
                      qtyMarkup += '<option value="' + j + '"' + selected + '>' + j + '</option>';
                    } else {
                      if(qtyMatched) {
                        qtyMarkup += '<option value="' + j + '"' + selected + '>' + j + '</option>';
                        break;
                      }
                    }
                }
                qtyMarkup += '</select>';
                markup +=   '<div class="left center" style="width: 125px;margin-top: 15px;">'
                        +       qtyMarkup
                        +   '</div>'
                        +   '<div class="left center" style="width: 160px;margin-top: 15px;">'
                        +       utilObj.formatMoney(subTotal, true)
                        +   '</div>'
                        + '</div>';  
                total += subTotal;
            }
            markup += '<div class="bag_total upper">'
                    +   '<span style="margin-right: 240px; margin-left: 570px;">Total</span>'
                    +   utilObj.formatMoney(total, true)
                    + '</div>';  
        }
        $('.bag_items').html(markup);
        $('.remove_item').unbind('click').click(function() {
            var prodIndex = $(this).data('index');
            cartObj.removeFromCart(prodIndex);
        });
        $('select.cart_qty').selectric();
        $('select.cart_qty').change(function(){
            cartObj.updateQuantity($(this).data('index'), $(this).val());
        });
    };
};

Checkout = function(addressesStr) {
    var PAYMENT_CASH = 0, PAYMENT_CREDIT_CARD = 1, 
            PAYMENT_DEBIT_CARD = 2, PAYMENT_NETBANKING = 3;
    var orderPlaced = false;
    var addresses = null;
    var selectedAddressId = 0;
    if (addressesStr != "") {
        addresses = JSON.parse(addressesStr);
    } 

    this.showAddressChoices = function() {
        var markup = '';
        for(var i=0; i < addresses.length; i++) {
          if(addresses[i].Id == 0) {
            break;
          }
          var checked = '';
          if (i == 0) {
              checked = 'checked="checked"';
          }
          markup += '<div id="existing_address_' + addresses[i].Id + '" class="relative">';
          markup += '<div style="position: absolute;">'
                 + '<input type="radio" value="' + addresses[i].Id + '" name="existingAddress" data-name="' 
                 + addresses[i].Name + '" + data-address="' + addresses[i].Address + '" data-city="' 
                 + addresses[i].City + '" data-state="' + addresses[i].State + '" data-pincode="'
                 + addresses[i].Pincode + '" data-mobile="' + addresses[i].Mobile + '"' 
                 + checked + '>'
                 + '</div>';
          markup += '<div style="font-size: 14px; margin-left: 40px;">'
                 + addresses[i].Name + ' - ' + addresses[i].Mobile + '<div>'
                 + addresses[i].Address + ', ' + addresses[i].City + ', ' + addresses[i].State
                 + ' - ' + addresses[i].Pincode
                 + '</div>'
                 + '</div>';
          markup += '<div class="remove_address like_link" data-addressid="' + addresses[i].Id + '">remove</div>';
          markup += '</div>';
        }
        $('.existing_addresses').html(markup);
        $('.remove_address').unbind('click').click(function() {
            var addressId = $(this).data('addressid');
            userObj.removeAddress( addressId );
            $('#existing_address_' + addressId).fadeOut();
        });
        $("input[type='radio'][name='existingAddress']").unbind('change').change(function() {
           checkoutObj.populateAddressInfo();
        });
        this.populateAddressInfo();
    };
    this.populateAddressInfo = function() {
        var selected = $("input[type='radio'][name='existingAddress']:checked");
        if (selected.length > 0) {
            $('#checkout_username').val( selected.data('name') );
            $('#checkout_address').val( selected.data('address') );
            $('#checkout_city').val( selected.data('city') );
            $('select.address_options_state').val( selected.data('state') );
            $('#checkout_pincode').val( selected.data('pincode') );
            $('#checkout_mobile').val( selected.data('mobile') );
            selectedAddressId = selected.val();
        }
    };
    this.showInitialState = function() {
        if(userObj.isLoggedIn()) {
            var email = userObj.getEmail();
            $('#checkout_email').val(email);
            $('.checkout_summary_email').text(email);
            $('.checkout_change_email').show();
            $('.checkout_info_address').show();

            if(addresses != null) {
              this.showAddressChoices();
            }
        } else {
            $('.checkout_info_email').show();
        }
        $('.checkout_change_email').unbind('click').click(function() {
            $('.checkout_req_info').not('.checkout_info_email').slideUp();
            $('.checkout_info_email').slideDown();
        });
        $('.move_next_email').unbind('click').click(function() {
            var email = $('#checkout_email').val();
            if (email == "") {
                utilObj.showErrorMessage("Please enter a valid email address");
                return;
            }
            $('.checkout_info_email').slideUp();
            $('.checkout_change_email').show();
            $('.checkout_summary_email').text(email);
            $('.checkout_info_address').slideDown();
        });
        $('.checkout_change_address').unbind('click').click(function() {
            $('.checkout_req_info').not('.checkout_info_address').slideUp();
            $('.checkout_info_address').slideDown();
        });
        $('.move_next_address').unbind('click').click(function() {
            var name = $('#checkout_username').val();
            var address = $('#checkout_address').val();
            var city = $('#checkout_city').val();
            var state = $('select.address_options_state').val();
            var pincode = $('#checkout_pincode').val();
            var mobile = $('#checkout_mobile').val();

            if (name == "" || address == "" || city == "" || state == "" || pincode == "" || mobile == "") {
                utilObj.showErrorMessage("Please enter valid values for all address fields");
                return;
            }
            var addressSummary = name + ' - ' + mobile + '<div style="font-size: 14px;">'
                    + address + ', ' + city + ', ' + state + ' - ' + pincode + '</div>';
            $('.checkout_info_address').slideUp();
            $('.checkout_summary_address').html(addressSummary);
            $('.checkout_change_address').show();
            $('.checkout_info_order').slideDown();
            cartObj.showItemList(true, null);
        });
        $('.checkout_change_order').unbind('click').click(function() {
            $('.checkout_req_info').not('.checkout_info_order').slideUp();
            $('.checkout_info_order').slideDown();
        });
        $('.move_next_order').unbind('click').click(function() {
            $('.checkout_info_order').slideUp();
            $('.checkout_summary_order').text('');
            $('.checkout_change_order').show();
            $('.checkout_info_payment').slideDown();
        });
        $('.place_order').unbind('click').click(function() {
            checkoutObj.placeOrder();
        });
    };
    this.placeOrder = function() {
        if (orderPlaced) {
            return;
        }
        var email = $('#checkout_email').val();
        var name = $('#checkout_username').val();
        var address = $('#checkout_address').val();
        var city = $('#checkout_city').val();
        var state = $('select.address_options_state').val();
        var pincode = $('#checkout_pincode').val();
        var mobile = $('#checkout_mobile').val();
        var productsJson = cartObj.getProducts();
        var products = JSON.stringify(productsJson);
        var paymentMethod = PAYMENT_CASH;
        if(email === "" || name === "" || address === "" || city === ""||
            state === "" || pincode === "" || mobile === "" || products === "" || paymentMethod === "") {
            utilObj.showErrorMessage("Please verify your order details again. Some fields are missing.");
            return;
        }
        orderPlaced = true;
        if(selectedAddressId != 0) {
            for(var i=0; i < addresses.length; i++) {
                if(addresses[i].Id == selectedAddressId) {
                    if(email == addresses[i].Email && name == addresses[i].Name 
                            && address == addresses[i].Address && city == addresses[i].City
                            && state == addresses[i].State && pincode == addresses[i].Pincode
                            && mobile == addresses[i].Mobile) {
                        //All values matched. Go ahead and use address id.
                    } else {
                        //User modified some value. Store as a new address.
                        selectedAddressId = 0;
                    }
                    break;
                }
            }
        }
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/user/purchase',
            data: {'email': email, 'name': name, 'address': address,
                'city': city, 'state': state, 'pincode': pincode,
                'mobile': mobile, 'products': products,
                'paymentMethod': paymentMethod, 'addressId': selectedAddressId
                },
            success: this.afterPlacingOrder,
            dataType: 'json'
        });
    };
    this.afterPlacingOrder = function(result) {
        if(result.Success == false) {
            utilObj.showErrorMessage("Unable to place order. Please try again.");
            orderPlaced = false;
            return;
        }
        var order = JSON.parse(result.Data.order);
        if(order.PaymentMethod == 0) {
            var formMarkup = '<form id="purchase_success" action="' + configObj.BASE_URL + '/user/purchaseSuccess' + '" method="POST">'
                + '<input type="hidden" name="orderId" value="' + order.Id + '">'
                + '</form>'
            $('body').append(formMarkup);
            $('#purchase_success').submit();
        } else {
            console.log("Redirect to payzippy for payment");
        }
    };
};

Comment = function() {
    var commentsEndPoint = configObj.BASE_URL + '/comment';
    var showCommentsIn;
    var overwritePrev;

    this.INSPIRATION = 1;

    this.createCommentHtml = function(commentator, comment) {
        return '<div class="comment">'
            + '<span class="commentator">' + commentator
            + '</span>'
            + '<span class="comment_text" title="' + comment 
            + '">' + comment + '</span>'
            + '</div>';
    };

    this.afterPosting = function(result) {
        if (result.Success == false) {
            alert("Unable to post comment. Please refresh the page and try again");
            return;
        }
        var commentsHtml = commentObj.createCommentHtml(result.Data.commentator, 
            result.Data.comment);
        var commentsArea = document.getElementById(showCommentsIn);
        if(overwritePrev == true) {
            commentsArea.innerHTML =  commentsHtml;
        } else {
            commentsArea.innerHTML =  commentsArea.innerHTML + commentsHtml;
        }
        commentsArea.style.display = "block";
    };

    this.postComment = function(type, subjectId, comment, divForShowingComments, overwrite) {
        if(!userObj.isLoggedIn()) {
            userObj.showLoginPopup();
            return;
        }
        trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "comment");
        showCommentsIn = divForShowingComments;
        overwritePrev = overwrite;
        comment = profanityFilterObj.filterText(comment);
        $.ajax({
            type: "POST",
            url: commentsEndPoint + '/add',
            data: {'type': type, 'subjectId': subjectId, 'comment': comment},
            success: this.afterPosting,
            dataType: 'json'
        });
    };
    this.getComments = function(type, subjectId, divForShowingComments) {
        showCommentsIn = divForShowingComments;
        $.ajax({
            type: "POST",
            url: commentsEndPoint + '/get',
            data: {'type': type, 'subjectId': subjectId},
            success: this.afterFetching,
            dataType: 'json'
        });
    };
    this.afterFetching = function(result) {
        if (result.Success == false) {
            return;
        }
        var comments = JSON.parse(result.Data.comments);
        var commentsHtml = '';
        for (var i = 0; i < comments.length; i++) {
            if(comments[i].Comment == "") {
                break;
            }
            commentsHtml += commentObj.createCommentHtml(comments[i].UserName,
                    comments[i].Comment);
        }
        $('#' + showCommentsIn).html(commentsHtml);
    };
};

Config = function(static_url,base_url,ins_url,designs_url,current_page,fb_app_name) {
    this.STATIC_URL = static_url;
    this.BASE_URL = base_url;
    this.INSPIRATIONS_URL = ins_url;
    this.DESIGNS_URL = designs_url;
    this.CUR_PAGE = current_page;
    this.FB_APP_NAME = fb_app_name;
    
    this.GOOGLE = 1;
    this.YAHOO = 3;
};

Contest = function() {
    this.submitPrice = function() {
        var price = $('#guessed_pendant_price').val();
        if(price == "") {
            return;
        }
        var priceAlreadyGuessed = localStorage.getItem("pendant_price_guess");
        if(priceAlreadyGuessed !== null) {
            return;
        }
        localStorage.setItem("pendant_price_guess", price);
        trackObj.count("action", configObj.CUR_PAGE, "quote", "submit", price);
        if(userObj.getUserId() != 0) {
            this.updateContestDetails();
        } else {
            userObj.askEmail();
        }
    };
    this.likePage = function(liked) {
        localStorage.setItem("fb_like", liked);
        this.updateContestDetails();
    };
    this.sharePage = function() {
        localStorage.setItem("fb_share", 1);
        contestObj.updateContestDetails();
    };
    this.updateContestDetails = function() {
        var guessedPrice = localStorage.getItem("pendant_price_guess");
        if(guessedPrice === null) {
            $('.contest_price_guess_0').show();
            return;
        } else {
            $('.contest_price_guess_0').fadeOut();
            var refCount = userObj.getReferralCount();
            var fbid = userObj.getFbid();
            var liked = localStorage.getItem('fb_like');
            if(fbid == "") {
                $('.contest_connected').fadeOut();
                $('.contest_not_connected').show();
            } else {
                $('.contest_not_connected').fadeOut();
                $('.contest_connected').show();
            }
            if(liked == null || liked == "" || liked == 0) {
                $('.contest_liked').fadeOut();
                $('.contest_not_liked').show();
            } else {
                $('.contest_not_liked').fadeOut();
                $('.contest_liked').show();
            }
            
            if(refCount < 3) {
                $('.contest_referred').fadeOut();
                $('.contest_not_referred').show();
            } else {
                $('.contest_not_referred').fadeOut();
                $('.contest_referred').show();
            }
            $('.contest_price_guess_1').show();
        }
        this.trackEntry(guessedPrice);
    };
   
    this.trackEntry = function(price) {
        var tracked = localStorage.getItem("contest_price_guess_tracked");
        if(tracked !== null) {
            return;
        }
        if(userObj.getUserId() != 0) {
            trackObj.count("contest", "pendant_price_guess", price);
            localStorage.setItem("contest_price_guess_tracked", true);
            var cogsURL = configObj.BASE_URL + "/contest/winPendant";
            facebookWrapperObj.publishCOGS('participated_in', 'quiz', cogsURL);
        }
    };
};

Email = function() {
    var tabForEmailSend = 'email';
    this.setMessage = function(message) {
        $('#invite_status_message').text(message);
    };
    this.sendEmail = function(tab, list) {
        var message = '';
        tabForEmailSend = tab;
        if (list == '') {
            message = 'Please enter valid email addresses';
            this.setMessage(message);
            return;
        }
        message = 'Sending Emails';
        this.setMessage(message);
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/user/sendReferralEmail',
            data: {emails: list},
            success: emailObj.emailSendSuccess,
            dataType: 'json'
        });
        var emailList = list.split(",");
        trackObj.count("action", configObj.CUR_PAGE, "invite", "email", tab, "", emailList.length - 1);
    };
    
    this.emailSendSuccess = function(result)
    {
        emailObj.setMessage(result.Data.message);
    };
};

Experiment = function(experimentsAccessed) {
    var expsAccessed = null;
    var syncId = 3;
    if (experimentsAccessed != "") {
        expsAccessed = JSON.parse(experimentsAccessed);
    }
    var baseKey = 'myExperiments';
    var allExpData = null;

    this.getBaseKey = function() {
        return baseKey;
    };

    this.getSyncUrl = function() {
        return "/experiment";
    };
    
    this.isTemp = function() {
        return false;
    };
    
    this.getAllExperiments = function() {
        if(allExpData != null) {
            return allExpData;
        }
        var allExperiments = '';
        var storageKey = baseKey;
        var exps = localStorage.getItem(storageKey);
        if (exps !== null) {
            allExperiments += exps;
        }

        storageKey = utilObj.getLocalKey(baseKey);
        exps = localStorage.getItem(storageKey);
        if (exps !== null) {
            allExperiments += exps;
        }

        storageKey = utilObj.getTempKey(baseKey);
        exps = localStorage.getItem(storageKey);
        if (exps !== null) {
            allExperiments += exps;
        }

        var allExps = allExperiments.split(";");
        allExpData = {};
        for (var i = 0; i < allExps.length; i++) {
            var expData = allExps[i].split(":");
            allExpData[ expData[0] ] = expData[1];
        } 
        return allExpData;
    };

    this.isPresent = function(exp) {
        var allExps = this.getAllExperiments();
        if(allExps[exp] !== undefined) {
            return true;
        }
        return false;
    };
    
    this.addToExperiments = function(exp, bucket) {
        var storageKey = utilObj.getLocalKey(baseKey);
        if (syncObj.isSyncInProgress(syncId)) {
            storageKey = utilObj.getTempKey(baseKey);
        }
        var curValue = localStorage.getItem(storageKey);
        if (curValue === null) {
            curValue = '';
        }
        curValue += exp.name + ':' + bucket + ';';
        localStorage.setItem(storageKey, curValue);
        allExpData[exp.name] = bucket;
        return true;
    };

    this.allocateBucket = function(exp) {
        var randomInt = utilObj.getRandomInt(1,100);
        var cumulativeWeight = 0;
        for (var bucket = 0; bucket < exp.variants; bucket++) {
            cumulativeWeight += exp.thresholds[bucket];
            if (randomInt < cumulativeWeight) {
                this.addToExperiments(exp, bucket);
                trackObj.count("experiment", exp.name, bucket);
                break;
            }
        }
    };
    
    this.getBucket = function(expName,expObj) {
        if(!this.isPresent(expName)) {
            this.allocateBucket(expObj);
        }
        return allExpData[expName];
    };
    
    this.processAllAccessed = function() {
        for(var expIndex in expsAccessed) {
            var exp = expsAccessed[expIndex];
            var bucket = this.getBucket(exp.name, exp);
            var experimentClass = '.jools_exp_' + exp.name + '_' + bucket;
            $(experimentClass).show();
        }
    };
};

Favorite = function() {
    var baseKey = 'myFavInspirations';
    var syncId = 0;

    this.getBaseKey = function() {
        return baseKey;
    };
    this.getSyncUrl = function() {
        return "/favorite";
    };
    this.isTemp = function() {
        return false;
    };
    this.MergeForSync = function(localData, tempData) {
       return localData + tempData;
    };

    this.getAllFavorites = function(hashed) {
        var allFavorites = '';
        var storageKey = baseKey;
        var favs = localStorage.getItem(storageKey);
        if (favs !== null) {
            allFavorites += favs;
        }

        storageKey = utilObj.getLocalKey(baseKey);
        favs = localStorage.getItem(storageKey);
        if (favs !== null) {
            allFavorites += favs;
        }

        storageKey = utilObj.getTempKey(baseKey);
        favs = localStorage.getItem(storageKey);
        if (favs !== null) {
            allFavorites += favs;
        }

        var allFavs = allFavorites.split(";");
        if (hashed == true) {
            var hashedFavs = {};
            for (var i = 0; i < allFavs.length; i++) {
                hashedFavs[allFavs[i]] = 1;
            }
            return hashedFavs;
        }
        return allFavs;
    };

    this.getRemovedFavorites = function() {
        var removedFavs = '';
        var storageKey = utilObj.getRemoveKey(baseKey);
        var removed = localStorage.getItem(storageKey);
        if (removed !== null) {
            removedFavs += removed;
        }

        return removedFavs.split(';');
    };

    this.isFavorited = function(inspirationId) {
        var allFavs = this.getAllFavorites(false);
        for (var i = 0; i < allFavs.length; i++) {
            if (inspirationId == allFavs[i]) {
                return true;
            }
        }
        return false;
    };

    this.addToFavorites = function(inspirationId) {
        if (inspirationId == undefined) {
            console.log("Inspiration id is not defined");
            return false;
        }
        var allFavorites = this.getAllFavorites(false);
        for (var i = 0; i < allFavorites.length; i++) {
            if (inspirationId == allFavorites[i]) {
                console.log("Attempt to add duplicate inspiration id", inspirationId);
                return false;
            }
        }

        var storageKey = utilObj.getLocalKey(baseKey);
        if (syncObj.isSyncInProgress(syncId)) {
            storageKey = utilObj.getTempKey(baseKey);
        }
        var curValue = localStorage.getItem(storageKey);
        if (curValue === null) {
            curValue = '';
        }
        curValue += inspirationId + ';';
        localStorage.setItem(storageKey, curValue);

        storageKey = utilObj.getRemoveKey(baseKey);
        this.checkAndRemove(storageKey, inspirationId);
        return true;
    };

    this.removeFromFavorites = function(inspirationId) {
        if (inspirationId == undefined) {
            console.log("Inspiration id is not defined");
            return false;
        }
        var removedFavs = this.getRemovedFavorites();
        for (var i = 0; i < removedFavs.length; i++) {
            if (inspirationId == removedFavs[i]) {
                console.log("Attempt to remove again", inspirationId);
                return false;
            }
        }

        var storageKey = utilObj.getRemoveKey(baseKey);
        var removedList = localStorage.getItem(storageKey);
        if (removedList === null) {
            removedList = '';
        }
        localStorage.setItem(storageKey, removedList + inspirationId + ';');

        storageKey = utilObj.getTempKey(baseKey);
        var matchFound = this.checkAndRemove(storageKey, inspirationId);
        if (matchFound) {
            return true;
        }

        storageKey = utilObj.getLocalKey(baseKey);
        matchFound = this.checkAndRemove(storageKey, inspirationId);
        if (matchFound) {
            return true;
        }

        storageKey = baseKey;
        matchFound = this.checkAndRemove(storageKey, inspirationId);
        return true;
    };

    this.checkAndRemove = function(storageKey, inspirationId) {
        var curList = localStorage.getItem(storageKey);
        var newList = '';
        var matchFound = false;
        if (curList !== null) {
            var temp = curList.split(';');
            for (var i = 0; i < temp.length; i++) {
                if (temp[i] == inspirationId) {
                    matchFound = true;
                } else if (temp[i] != '') {
                    newList += temp[i] + ';';
                }
            }
        }
        if (matchFound) {
            localStorage.setItem(storageKey, newList);
        }
        return matchFound;
    }
};

FacebookWrapper = function() {
	this.GRAPH_API_BASE = "https://graph.facebook.com";
    var friendsData = null;
    var allContactsHtml = null;
    var totalFriendCount = 0;
    var invitesJustSentTo = new Array();
    this.numFbInviteDialogs = 0;
    this.totalInvitesSent = 0;
    var contactsLimit = 100;
    if (!utilObj.isMobileBrowser()) {
        contactsLimit = 1000;
    }
    this.tryFBLogin = function(clickSrc, callback) {
        var userId = userObj.getUserId();
        var signState = "sign_in";
        if(callback == undefined) {
            callback = facebookWrapperObj.afterFBLogin;
        }
        if(userId == 0) {
            signState = "sign_up";
            trackObj.count("action", configObj.CUR_PAGE, "signup_intent", clickSrc, "fb_click");
        }
        trackObj.count("action", configObj.CUR_PAGE, signState, clickSrc, "fb_click");
        FB.login(function(response) {
            var userId = userObj.getUserId();
            var signState = "sign_in";
            if(userId == 0) {
                signState = "sign_up";
            }
            if (response.authResponse) {
                trackObj.count("action", configObj.CUR_PAGE, signState, clickSrc, "fb_gave_perms");
                $.ajax({
                    type: "POST",
                    url: configObj.BASE_URL + "/user/signInWithFB",
                    data: {
                        'userID': response.authResponse.userID, 
                        'jref': userObj.referrer,
                        'signedRequest': response.authResponse.signedRequest,
                        'accessToken': response.authResponse.accessToken
                    },
                    success: callback,
                    dataType: 'json'
                });
            } else {
                trackObj.count("action", configObj.CUR_PAGE, signState, clickSrc, "fb_didnot_give_perms");
            }
        }, {scope: 'email,publish_actions'});
    };
    this.afterFBLogin = function(res) {
        if(res.Success == true) {
            if (utilObj.isMobileBrowser()) {
                window.location = configObj.BASE_URL + '/';
            } else {
                utilObj.reloadPage();
            }
        } else {
            alert("Unable to login, please try again: " + res.Data.message);
        }
    };
    this.inviteFriends = function(src) {
        if(src == undefined) {
            src = configObj.CUR_PAGE;
        }
        trackObj.count("action", src, "invite", "fb_request", "click");
        this.getFBPermissions(this.showInviteFBFriends);
    };
    this.inviteFriendsDefaultFB = function() {
        FB.ui({method: 'apprequests',
            message: 'Come and check out Jools'
            }, this.inviteFriendsCallback);
    };
    this.inviteFriendsCallback = function(result) {
        facebookWrapperObj.numFbInviteDialogs--;
        if(result.to == undefined) {
            trackObj.count("action", configObj.CUR_PAGE, "invite", "fb_request", "abort");
            return;
        }
        var numInvites = result.to.length;
        facebookWrapperObj.totalInvitesSent += numInvites;
        trackObj.count("action", configObj.CUR_PAGE, "invite", "fb_request", "success", "", numInvites);
        var requestId = result.request;
        userObj.storeFBRequest(requestId, result.to);
        if(facebookWrapperObj.numFbInviteDialogs <= 0) {
            userObj.invitesSent(facebookWrapperObj.totalInvitesSent);
        }
    };
    this.sortFBFriendsByName = function(a,b) {
        var x = a.name.toLowerCase();
        var y = b.name.toLowerCase();
        return ((x < y) ? -1 : ((x > y) ? 1 : 0));
    };
    this.getFriends = function() {
        if (friendsData != null) {
            facebookWrapperObj.showInviteFBFriends();
            return;
        }
        FB.api('/me/friends', function(response) {
            if (response.data) {
                friendsData = response.data.sort(facebookWrapperObj.sortFBFriendsByName);
                facebookWrapperObj.showInviteFBFriends();
            }
        });
    };
    this.showInviteFBFriends = function() {
        if (utilObj.isMobileBrowser()) {
            if (friendsData == null) {
                window.location = configObj.BASE_URL + '/user/inviteFBFriends?utm_source=' + configObj.CUR_PAGE;
            } else {
                facebookWrapperObj.populateFBFriends();
            }
        } else {
            if (friendsData == null) {
                facebookWrapperObj.getFriends();
                return;
            }
            utilObj.createPopup(configObj.BASE_URL + "/user/inviteFBFriends",
                    {}, 900, 600, true);
        }
    };
    this.populateFBFriends = function() {
        if (friendsData == null) {
            this.getFriends();
            return;
        }
        totalFriendCount = 0;
        allContactsHtml = '';
        $.each(friendsData, function(index, friend) {
            if (totalFriendCount < contactsLimit) {
                totalFriendCount++;
                allContactsHtml += facebookWrapperObj.getContactHtml(friend);
            }
        });
        $('#fb_friends').html(allContactsHtml);
        $('#fb_total_friend_count').text(totalFriendCount);
        this.bindInviteListeners();
        this.selectAllFB();
    };
    this.getContactHtml = function(friend) {
        var markup = ''
                + '<div class="fb_contact_info fb_contact_' + friend.id + '">'
                + '<img class="fb_contact_image" src="https://graph.facebook.com/'
                + friend.id + '/picture" height="60" width="60">'
                + '<div class="fb_checkbox right">'
                + '<input type="checkbox" value="' + friend.id + '" name="fb_contact">'
                + '</div>'
                + '<div class="fb_contacts_selector_name" title="' + friend.name + '">'
                + friend.name
                + '</div>'
                + '</div>';
        return markup;	
    };
    this.showMatchedFriendList = function(partialName) {
        partialName = partialName.toLowerCase();
        $.each(friendsData, function(index, friend) {
            var name = friend.name.toLowerCase();
            var n = name.search(partialName);
            if (n != -1) {
                $('.fb_contact_' + friend.id).show();
            } else {
                $('.fb_contact_' + friend.id).hide();
            }
        });
    };
    this.showFullFriendList = function() {
        $('.fb_contact_info').show();
    };
    this.selectAllFB = function() {
        var count = 0;
        $("input:checkbox[name=fb_contact]").each(function() {
            $(this).prop('checked', true);
            count++;
        });
        $('#fb_contacts_selected_count').text(count);
        $('#fb_invite_progress_bar_bg').attr('style', 'width: 100%;');
    };
    this.unselectAllFB = function() {
        $("input:checkbox[name=fb_contact]").prop('checked', false);
        $('#fb_contacts_selected_count').text(0);
        $('#fb_invite_progress_bar_bg').attr('style', 'width: 1%;');
    };
    this.bindInviteListeners = function() {
        $('#fb_contacts_master_checkbox').unbind('change').change(function() {
            if ($(this).is(":checked")) {
                facebookWrapperObj.selectAllFB();
            } else {
                facebookWrapperObj.unselectAllFB();
            }
        });
        $("input:checkbox[name=fb_contact]").unbind('change').change(function() {
            var count = parseInt($('#fb_contacts_selected_count').text());
            if ($(this).is(":checked")) {
                count++;
            } else {
                count--;
            }
            if (count < 0) {
                count = 0;
            }
            $('#fb_contacts_selected_count').text(count);
            var totalCount = parseInt($('#fb_total_friend_count').text())
            var width = (count / totalCount) * 100;
            if (width < 95) {
                width += 1;
            }
            $('#fb_invite_progress_bar_bg').attr('style', 'width: ' + width + '%;');
        });
    };
    this.sendFBInvites = function() {
        var curFbidList = '';
        var selectedFriendCount = 0;
        facebookWrapperObj.numFbInviteDialogs = 0;
        facebookWrapperObj.totalInvitesSent = 0;
        $("input:checkbox[name=fb_contact]:checked").each(function() {
            var curFbid =  $(this).val();
            if(invitesJustSentTo[curFbid] == undefined) {
                invitesJustSentTo[curFbid] = 1;
                selectedFriendCount++;
                if (selectedFriendCount % 50 == 0) {
                    curFbidList +=  curFbid;
                    facebookWrapperObj.numFbInviteDialogs++;
                    FB.ui({
                        method: 'apprequests',
                        message: 'Come and check out Jools',
                        to: curFbidList
                    }, facebookWrapperObj.inviteFriendsCallback);
                    curFbidList = '';
                } else {
                    curFbidList +=  curFbid + ',';
                }
            }
        });
        if(selectedFriendCount % 50 != 0) {
            facebookWrapperObj.numFbInviteDialogs++;
            FB.ui({
                method: 'apprequests',
                message: 'Come and check out Jools',
                to: curFbidList
            }, facebookWrapperObj.inviteFriendsCallback);
        }
        trackObj.count("action", "invite_popup", "invite", "fb_request_start", "click");
    };
    this.getFBPermissions = function(callback) {
        if(userObj.getFbid() == "") {
            facebookWrapperObj.tryFBLogin('connect', callback);
            return;
        }
        FB.getLoginStatus(function(response) {
            if(response.status != 'connected') {
                facebookWrapperObj.tryFBLogin('invite', callback);
            } else {
                if(callback !== null) {
                    callback();
                }
            }
        });
    };
    this.publishCOGS = function(action, obj, objUrl) {
        var cogsUrl = this.GRAPH_API_BASE + '/me/' + configObj.FB_APP_NAME + ':' + action;
        var params = JSON.parse('{"' + obj + '": "' + objUrl + '"}');
        if (typeof FB == "undefined" || FB.getAccessToken() == null) {
            setTimeout(function() {
                facebookWrapperObj.publishCOGS(action, obj, objUrl);
            }, 1000);
            return;
        }
        FB.api(cogsUrl, 'post', params, function(response) {
            if (!response) {
                trackObj.count("og", action, obj, "fail","unknown");
            } else if (response.error) {
                trackObj.count("og", action, obj, "fail", response.error.code);
            } else {
                trackObj.count("og", action, obj, "success");
            }
        });
    };
};

Inspiration = function(inspirations) {
    var initialInspirations = null;
    if(inspirations != "") {
       initialInspirations = JSON.parse(inspirations);
    }
    
    var currentWidth = 0;
    var inspirationsBaseUrl = configObj.BASE_URL + "/inspiration";
    var inspirationsUrl = configObj.INSPIRATIONS_URL;
    this.allLotsDone = false;
    this.requestInProgress = false;
    var addToFav = "Add to favorites";
    var removeFromFav = "Remove from favorites";
    var numOfViewsText = "Number of views";
    var shareOnFacebook = "Share on Facebook";
    var placeholderComment = "What do you think of this piece?";
    var commentType = 1;
    var lastId = null;
    var crossPromoOptionCount = 3;
    var curCrossPromo = 0;
    var intervalBeforePromo = 10;
   
    this.getCrossPromoPin = function() {
        var pinHtml = '';
        if (curCrossPromo == 0) {
            pinHtml += ''
                    + '<div class="grid ins_cross_promo win_pendant">'
                    +  '<img src= "'+ configObj.STATIC_URL 
                    + '/images/home_win_pendant.jpg" alt="Win this pendant" style="max-width: 100%;">'
                    +  '<div class="cross_promo_wrapper">'
                    +    '<div class="goudy_text ins_cross_promo_head">'
                    +    'Win this pendant'
                    +    '</div>'
                    +    '<div class="lucida_text cross_promo_desc">'
                    +     'Play our fun guessing game'
                    +     '<div> & win this diamond pendant! </div>'
                    +    '</div>'
                    +    '<button class="win_pendant_link light_blue_button" data-src="ins_xpromo"'
                    +    'data-medium="' + configObj.CUR_PAGE  + '" style="width: 120px;">'
                    +     'Guess'
                    +    '</button>'
                    +   '</div>'
                    + '</div>';
        } else if (curCrossPromo == 1) {
            if (!utilObj.isMobileBrowser()) {
                pinHtml += ''
                    + '<div class="grid ins_cross_promo play_hot">'
                    +  '<img src= "'+ configObj.STATIC_URL 
                    + '/images/home_hot_not.jpg" alt="Hot or not" style="max-width: 100%;">'
                    +  '<div class="cross_promo_wrapper">'
                    +    '<div class="goudy_text ins_cross_promo_head">'
                    +    'Test your style quotient'
                    +    '</div>'
                    +    '<div class="lucida_text cross_promo_desc">'
                    +     'Rate jewellery images'
                    +     '<div> & unlock mystery gifts </div>'
                    +    '</div>'
                    +    '<button class="play_hot_link light_blue_button" data-src="ins_xpromo"'
                    +    'data-medium="' + configObj.CUR_PAGE  + '" style="width: 120px;">'
                    +     'Play'
                    +    '</button>'
                    +   '</div>'
                    + '</div>';
            }
        } else if (curCrossPromo == 2) {
            pinHtml += ''
                    + '<div class="grid ins_cross_promo earn_cash">'
                    +  '<img src= "'+ configObj.STATIC_URL 
                    + '/images/home_gift_500.jpg" alt="Gift 500,Get 500" style="max-width: 100%;">'
                    +  '<div class="cross_promo_wrapper">'
                    +    '<div class="goudy_text ins_cross_promo_head">'
                    +    'Gift <span class="lucida_text">500</span>, Get <span class="lucida_text">500</span>!'
                    +    '</div>'
                    +    '<div class="lucida_text cross_promo_desc">'
                    +     'Gift your friends <span class="WebRupee">Rs.</span>500 &'
                    +     '<div>earn up to <span class="WebRupee">Rs.</span>5000 in cash credits</div>'
                    +    '</div>'
                    +    '<button class="earn_cash_link light_blue_button" data-src="ins_xpromo"'
                    +    'data-medium="' + configObj.CUR_PAGE  + '" style="width: 120px;">'
                    +     'Gift'
                    +    '</button>'
                    +   '</div>'
                    + '</div>';
        }

        curCrossPromo = (curCrossPromo + 1) % crossPromoOptionCount;
        return pinHtml;
    };
    
    this.makeInspiration = function(data) {
        //var lastComment = '';
        //var commentStyle = 'background: white;'
        //var commentorId = '';
        //var commentorName = '';
        var uploaderInfo = '';
        if (data.EncId == "") {
            return '';
        }
        /*if (data.LastComment != '' && data.LastComment != null) {
            lastComment = data.LastComment;
            commentorId = data.LastCommentorId;
            commentorName = data.LastCommentorName;
            commentStyle = '';
        }*/
        if (data.UploaderName != "") {
            uploaderInfo = '<div class="meta">' + utilObj.getFirstName(data.UploaderName) + '</div>';
        }
        var favClass = 'fav_' + data.EncId;
        favClass = utilObj.getJqueryFriendlyEncId(favClass);
        var tooltip = "tooltip";
        if(utilObj.isMobileBrowser()) {
            tooltip = "";
        }
        var html = ''
                + '<div class="grid hide_inspiration">'
                +   '<div class="imgholder clickable" data-subjectid="'
                +       data.EncId + '">'
                +       '<img src="' + inspirationsUrl + '/small/' + data.ImageName + '" alt="">'
                +       '<div class="imgholder_hover">'
                +           '<div class="ins_quick_look">'
                +               'Quick look'
                +           '</div>'
                +       '</div>'
                +   '</div>'
                +   '<div class="ins_info">'
                +       uploaderInfo
                //+ '<div id="board_viewed_' + data.EncId
                //+ '" class="viewed tooltip" title="' + numOfViewsText
                //+ '" data-subjectid="' + data.EncId + '">'
                //+ data.ViewedCount
                //+ '</div>'
                +       '<div class="favorite clickable ' + tooltip + ' ' + favClass
                +           '" title="' + addToFav + '" '
                +           'data-subjectid="' + data.EncId + '">'
                +          data.FavoritedCount
                +       '</div>'
                +   '</div>'
                +   '<div class="like_option favorite clickable ' + tooltip + ' ' + favClass
                +           '" data-subjectid="' + data.EncId + '">'
                +       'Like'
                +   '</div>'
                +   '<div class="ins_sharing_options">'
                +       '<div class="fb_share_inspiration clickable" title="'
                +       shareOnFacebook + '" '
                +       'data-subjectid="' + data.EncId + '" data-picture="'
                +       data.ImageName + '">'
                +       '<div class="ins_share_text">Share</div>'
                +       '</div>'
                //+ '<div class="pin_share_inspiration clickable" title="Share on Pinterest" '
                //+ 'data-subjectid="' + data.EncId + '" data-picture="'
                //+ data.ImageName + '">'
                //+ '</div>'
                +   '</div>'
                //+ '<div class="comment" style="' + commentStyle + '" '
                //+ 'id="last_comment_' + data.EncId + '">'
                //+ '<span class="commentator">' + commentorName
                //+ '</span>'
                //+ '<span class="comment_text" title="' + lastComment + '">' 
                //+ lastComment + '</span>'
                //+ '</div>'
                //+ '<div style="margin-left: 5px;">'
                //+ '<input class="comment_box" type="text" '
                //+ 'style="width: 98%; font-size: 14px; font-weight: normal;" '
                //+ 'maxlength="1024" data-subjectid="' + data.EncId
                //+ '" data-type="' + commentType + '"'
                //+ 'placeholder="' + placeholderComment + '">'
                //+ '</div>'
                + '</div>';
        return html;
    };
    
    this.makeApprovalInspiration = function(data) {
        if (data.Id == "") {
            return '';
        }
        var html = ''
                + '<div class="grid">'
                + '<div class="imgholder">'
                + '<img src="' + configObj.BASE_URL + '/imageServer/uploads/' 
                + data.ImageName + '" alt="">'
                + '</div>'
                + '<button class="approveImage light_pink_button" data-id="' + data.Id 
                + '">Approve</button>'
                + '<button class="rejectImage light_pink_button"  data-id="' + data.Id
                + '">Reject</button>'
                + '</div>';
        return html;
    };
    
    this.makeTaggingInspiration = function(data) {
        if (data.EncId == "") {
            return '';
        }
        var encIdForJquery = utilObj.getJqueryFriendlyEncId(data.EncId);
        var categoryList = tagObj.getCategoryList();
        var checkboxes = '';
        for (var category in categoryList) {
            checkboxes += "<div style='margin-bottom: 20px;background: #AACDEC;'>" 
                    + category + ": ";
            for(var i=0; i < categoryList[category].length; i++) {
                checkboxes += '<input type="checkbox" name="' 
                        + category + '" value="' + categoryList[category][i].Tag + '"'
                        + 'class="' + encIdForJquery + '" data-tagid="' 
                        + categoryList[category][i].Id + '">'
                        + categoryList[category][i].Tag;
            }
            checkboxes += "</div>";
        }
        var html = ''
                + '<div class="grid">'
                + '<div class="imgholder">'
                + '<img src="' + configObj.BASE_URL + '/images/inspirations/' 
                + data.ImageName + '" alt="">'
                + '</div>'
                + checkboxes
                + '<button class="tagInspiration light_pink_button" data-id="' + data.EncId 
                + '">Tag</button>'
                + '</div>';
        return html;
    };

    this.showInspirationBoard = function() {
        if(initialInspirations == null) {
            return;
        }
        var boardHtml = '';
        for (var i = 0; i < initialInspirations.length; i++) {
            boardHtml += this.makeInspiration( initialInspirations[i] );
            lastId = initialInspirations[i].EncId;
            if (i > 0 && i % intervalBeforePromo == 0) {
                boardHtml += this.getCrossPromoPin();
            }
        }
        $('#allInspirations').html(boardHtml);
        this.setupGrid();
        $(window).bind('scroll', this.onScroll);
    };
    
    this.showInspirationApprovalBoard = function() {
        if (initialInspirations == null) {
            return;
        }
        var boardHtml = '';
        for (var i = 0; i < initialInspirations.length; i++) {
            if(initialInspirations[i].Id == 0) {
                break;
            }
            boardHtml += this.makeApprovalInspiration(initialInspirations[i]);
        }
        $('#allInspirations').html(boardHtml);
        this.setupGrid();
    };
    
    this.showInspirationTaggingBoard = function() {
        if (initialInspirations == null) {
            return;
        }
        var boardHtml = '';
        for (var i = 0; i < initialInspirations.length; i++) {
            if(initialInspirations[i].EncId == '') {
                break;
            }
            boardHtml += this.makeTaggingInspiration(initialInspirations[i]);
        }
        $('#allInspirations').html(boardHtml);
        this.setupGrid();
    };

    this.markApproved = function(approve) {
        approve.parent().hide();
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + "/webadmin/markApproved",
            data: {'id': approve.data('id')}
        });
    };
    
    this.markRejected = function(reject) {
        reject.parent().hide();
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + "/webadmin/markRejected",
            data: {'id': reject.data('id')}
        });
    };
    
    this.tagInspiration = function(tagButton) {
        var inspirationId = tagButton.data('id');
        var jqueryFriendlyId = utilObj.getJqueryFriendlyEncId(inspirationId);
        var tags = '';
        var tagIds = '';
        var tagCategoryCount = 0;
        var categoriesTagged = new Array();
        $('.' + jqueryFriendlyId).each(function() {
            if($(this).is(":checked")) {
                if(categoriesTagged[ $(this).attr('name') ] == undefined) {
                    tagCategoryCount++;
                    categoriesTagged[ $(this).attr('name') ] = true;
                }
                tags += $(this).val() + ",";
                tagIds += $(this).data('tagid') + ',';
            }
        });
        if(tagCategoryCount < 4) {
            alert("You must choose atleast 1 tag from each category");
            return;
        }
        tagButton.parent().hide();
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + "/webadmin/tagInspiration",
            data: {'encId': inspirationId, 'tags': tags, 'tagIds': tagIds}
        });
    };

    this.setupGrid = function() {
        //blocksit define
        $('#allInspirations').BlocksIt({
            numOfCol: 3,
            offsetX: 8,
            offsetY: 8
        });

        //window resize
        $(window).resize(function() {
            inspirationObj.handleResize();
        });
        this.handleResize();
        this.prepareCommentBoxes();
    };
    
    this.onScroll = function(event) {
        if (inspirationObj.allLotsDone || inspirationObj.requestInProgress) {
            return;
        }
        // Check if we're within x pixels of the bottom edge of the browser window.
        var winHeight = window.innerHeight ? window.innerHeight : $(window).height(); // iphone fix
        var closeToBottom = ($(window).scrollTop() + winHeight > $(document).height() - 800);

        if (closeToBottom) {
            inspirationObj.requestInProgress = true;
            $('.loading').show();
            trackObj.count("action", configObj.CUR_PAGE, "page_scroll");
            //Get the next items by making ajax call
            $.ajax({
                type: "POST",
                url: inspirationsBaseUrl + "/getNextLot",
                data: {'lastId': lastId},
                success: inspirationObj.afterNextLot,
                dataType: 'json'
            });
        }
    };
    
    this.afterNextLot = function(result) {
        if (result.Success == false) {
            inspirationObj.allLotsDone = true;
            return;
        }
        var nextInspirations = JSON.parse(result.Data.inspirations);
        var numInspirations = nextInspirations.length;
        var index = 0;
        var boardHtml = '';
        for (; index < numInspirations; index++) {
            boardHtml += inspirationObj.makeInspiration(nextInspirations[index]);
            lastId = nextInspirations[index].EncId;
            if (index > 0 && index % intervalBeforePromo == 0) {
                boardHtml += inspirationObj.getCrossPromoPin();
            }
        }
        $('#allInspirations').append(boardHtml);
        inspirationObj.prepareCommentBoxes();

        //currentWidth = 0;
        inspirationObj.handleResize();
        inspirationObj.requestInProgress = false;
        inspirationObj.markFavorited();
    };
    
    this.prepareCommentBoxes = function() {
        $('.comment_box').unbind("keydown").keydown(function(e) {
            if (e.keyCode == 13) {
                var comment = $(this).val();
                if (comment == '') {
                    return;
                }
                var type = $(this).data('type');
                var subjectId = $(this).data('subjectid');
                $(this).val('');
                commentObj.postComment(type, subjectId, comment,
                        'last_comment_' + subjectId, true);
            }
        });
    };
    
    this.handleFavoriting = function(subject) {
        if (!userObj.isLoggedIn()) {
            userObj.showLoginPopup();
            return;
        }
        var subjectId = subject.data('subjectid');
        var success = favoriteObj.addToFavorites(subjectId);
        if (!success) {
            return;
        }
        trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "favorite");
        var favClass = '.fav_' + subjectId;
        favClass = utilObj.getJqueryFriendlyEncId(favClass);
        $(favClass).each(function() {
            var count = $(this).text();
            count = parseInt(count);
            if(!isNaN(count)) {
                $(this).text((count + 1));
            }
            if($(this).hasClass('like_option')) {
                $(this).text('Liked');
            }
            $(this).addClass('favorited').removeClass('favorite')
                .tooltip("option", "content", removeFromFav).unbind('click')
                .click(function() {
                    inspirationObj.handleUnfavoriting($(this));
                });
        });
        var insUrl = configObj.BASE_URL + '/inspiration/hotTrends?id=' + subjectId;
        facebookWrapperObj.publishCOGS('loved', 'product', insUrl);
    };

    this.handleUnfavoriting = function(subject) {
        if (!userObj.isLoggedIn()) {
            alert("Please login to access this functionality");
            return;
        }
        var subjectId = subject.data('subjectid');
        var success = favoriteObj.removeFromFavorites(subjectId);
        if (!success) {
            return;
        }
        trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "unfavorite");
        var favClass = '.fav_' + subjectId;
        favClass = favClass.substring(0, favClass.length - 1);
        $(favClass).each(function() {
            var count = $(this).text();
            count = parseInt(count);
            if (!isNaN(count)) {
                $(this).text((count - 1));
            }
            $(this).addClass('favorite').removeClass('favorited')
                .tooltip("option", "content", addToFav).unbind('click')
                .click(function() {
                    inspirationObj.handleFavoriting($(this));
                });
            if($(this).hasClass('like_option')) {
                $(this).text('Like');
            }
         });
    };

    this.markFavorited = function() {
        if (!userObj.isLoggedIn()) {
            return false;
        }
        var hashedFavs = favoriteObj.getAllFavorites(true);
        $('.favorite').each(function() {
            var subjectId = $(this).data('subjectid');
            if (hashedFavs[subjectId] == 1) {
                $(this).addClass('favorited').removeClass('favorite')
                        .tooltip("option", "content", removeFromFav);
                var count = $(this).text();
                if (count !== "" && count === "0") {
                    $(this).text('1');
                }
                if ($(this).hasClass('like_option')) {
                    $(this).text('Liked');
                }
            }
        });
    };
    this.showInspirationInfo = function(subjectId) {
		utilObj.createPopup(inspirationsBaseUrl + "/getInitialInfo",
			{'subjectId': subjectId}, 900, 600, true);
	};
    this.handleResize = function() {
        $('.imgholder').unbind("click").click(function(e) {
            var subjectId = $(this).data('subjectid');
            var success = viewsObj.addToViews(subjectId);
            /*if (success == true) {
                var countElement = document.getElementById('board_viewed_' + subjectId);
                var count = countElement.innerHTML;
                countElement.innerHTML = parseInt(count) + 1;
            }*/
            trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "click");
			inspirationObj.showInspirationInfo(subjectId);
        });
	if (!utilObj.isMobileBrowser()) {
            $('.imgholder').unbind("hover").hover(function() {
                $(this).find('.imgholder_hover').show();
            }, function() {
                $(this).find('.imgholder_hover').hide();
            });
        }
        $('.fb_share_inspiration').unbind("click").click(function() {
            var subjectId = $(this).data('subjectid');
            trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "share_fb", "click");
            var pictureUrl = configObj.STATIC_URL  + "/images/inspirations/" + $(this).data('picture');
            var feedDescription = "I was going through the inspiration board at Jools, and found this incredible jewellery design! Check out more such designs at www.jools.in, or share your own designs and make some cash!";
            var feedTitle = "Look what I found on Jools!";
            FB.ui({
                method: 'feed',
                link: configObj.BASE_URL + "/inspiration/hotTrends?utm_source=fb_feed&id=" + subjectId,
                picture: pictureUrl,
                description: feedDescription,
                name: feedTitle,
            }, function(response){
                if(response == null) {
                    trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "share_fb", "abort");
                    return;
                }
		trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "share_fb", "success");
            });
        });
        $('.pin_share_inspiration').unbind('click').click(function() {
            var subjectId = $(this).data('subjectid');
            trackObj.count("action", configObj.CUR_PAGE, "jool", subjectId, "share_pinterest");
            var pictureUrl = configObj.STATIC_URL  + "/images/inspirations/" + $(this).data('picture');
            var url = encodeURIComponent(configObj.BASE_URL);
            var description = "Check out this amazing jewellery design on www.jools.in"
            var buttonUrl = "https://www.pinterest.com/pin/create/button/?url=" + url
                + "&media=" + encodeURIComponent(pictureUrl) 
                + "&description=" + encodeURIComponent(description);
            window.open(buttonUrl, "Share on Pinterest", "height=300,width=800");
        });
        
        $('.favorite').unbind('click').click(function() {
            inspirationObj.handleFavoriting($(this));
        });

        $('.favorited').unbind('click').click(function() {
            inspirationObj.handleUnfavoriting($(this));
        });
        utilObj.bindLinks();
        utilObj.styleInputsAndButtons();
        $('#allInspirations').imagesLoaded(function() {
            inspirationObj.requestInProgress = false;
            $('.loading').hide();
            var winWidth = 1000;
            var conWidth = 1000;
            col = 3;
            
            if(utilObj.isMobileBrowser()) {
                //winWidth = window.innerWidth;
                winWidth = $('body').width();
                conWidth = '100%';
                if (winWidth < 320) {
                    col = 1;
                } else if (winWidth < 480) {
                    col = 1;
                } else if (winWidth < 640) {
                    col = 2;
                } else if (winWidth < 768) {
                    col = 2;
                } else if (winWidth < 1024) {
                    col = 3;
                } else {
                    col = 3;
                }
            }

            //if (conWidth != currentWidth) 
            {
                currentWidth = conWidth;
                $('#allInspirations').width(conWidth);
                $('#allInspirations').BlocksIt({
                    numOfCol: col,
                    offsetX: 8,
                    offsetY: 8
                });
            }
            $('.hide_inspiration').removeClass('hide_inspiration');
        });
    };
};

InviteTier = function() {
	/*var maxInvitesForCredit = 20;
	var tiers = [0,1,2,3,4];
	var requirements = [0,3,5,10,maxInvitesForCredit]; 
	var credits = [500,1000,1500,2500,5000]; 

	this.getTierForReferrals = function(referralCount) {
            var curTier = 0;
            for (var i = 0; i < tiers.length; i++) {
                if (referralCount >= requirements[i]) {
                    curTier = tiers[i];
                } else {
                    break;
                }
            }
            return curTier;
	};*/

	this.getProgressBarWidth = function(referralCount) {
            /*var tier = this.getTierForReferrals(referralCount);
            var tierWidth = 100.0 / (tiers.length);
            var pWidth = 0;
            if (userObj.getUserId() == 0) {
                pWidth = 0;
            } else if (referralCount >= maxInvitesForCredit) {
                pWidth = 100;
            } else {
                pWidth = tierWidth * (tier + 1) +
                    ((referralCount - requirements[tier]) * tierWidth) /
                    (requirements[tier + 1] - requirements[tier]);
            }*/
            var pWidth = 0;
            if(referralCount >= 9) {
                pWidth = 100;
            } else {
                pWidth = (referralCount + 1) / 10 * 100;
            }
            return pWidth;
	};
        this.getCashEarned = function() {
            var userId = userObj.getUserId();
            if(userId == 0) {
                return 0;
            }   
            return 500 * (userObj.getReferralCount() + 1);
        };
        this.getPossibleCash = function() {
            return 5000 - this.getCashEarned();
        };
};

Product = function(info) {
    var productInfo = info;
    var originalProductInfo = null;
    var selectedProductInfo = null;
    var curProductInfo = null;
    var productVariants = null;
    var curIndex = null;
    var curVariant = null;
    var curView = 1;
    var NO_STONE = -1;
    var thumbnailTimer = null;

    originalProductInfo = productInfo.ProductData;
    selectedProductInfo = jQuery.extend(true, {}, originalProductInfo);
    curProductInfo = jQuery.extend(true, {}, originalProductInfo);
    productVariants = productInfo.ProductVariants;
   
    var METAL_WHITE = 0, METAL_YELLOW = 1, METAL_ROSE = 2;
    var GOLD_22KT = 0, GOLD_18KT = 1, GOLD_14KT = 2, SILVER = 3;
    var metalColorList = {
        0: {Id: METAL_WHITE, Name: "White", ShortName: "W"},
        1: {Id: METAL_YELLOW, Name: "Yellow", ShortName: "Y"},
        2: {Id: METAL_ROSE, Name: "Rose", ShortName: "P"}
    };
    var metalList = {
        0: {Id: GOLD_22KT, Name: "Gold", ShortName: "G22", Karat: "22K"},
        1: {Id: GOLD_18KT, Name: "Gold", ShortName: "G18", Karat: "18K"},
        2: {Id: GOLD_14KT, Name: "Gold", ShortName: "G14", Karat: "14K"},
        3: {Id: SILVER, Name: "Silver", ShortName: "SLV", Karat: ""},
    };

    var DIAMOND = 0, BLACK_DIAMOND = 1, EMERALD = 2, RUBY = 3, SAPPHIRE_PINK = 4, SAPPHIRE_YELLOW = 5, SAPPHIRE_BLUE = 6;
    var CENTER_STONE = 0, ACCENT_STONE_ONE = 1, ACCENT_STONE_TWO = 2;
    var GEMSTONE_CAT_ZERO = 0, GEMSTONE_CAT_ONE = 1, GEMSTONE_CAT_TWO = 2;
    var GEM_CLARITY_ZERO = 0, GEM_CLARITY_ONE = 1, GEM_CLARITY_TWO = 2;
    var GEM_ROUND = 0, GEM_SQUARE = 1, GEM_MARQUISE =2, GEM_CUSHION =3, GEM_TRILLION = 4, GEM_BAGUTTE = 5;
    var GEM_SET_PRONG = 0, GEM_SET_BEZEL = 1,GEM_SET_PAVE = 2,GEM_SET_CHANNEL = 3;
    var gemstoneList = {
        0: {Id: DIAMOND, Name: "Diamond", ShortName: "DW", Birthstone: "April"},
        1: {Id: BLACK_DIAMOND, Name: "Black Diamond", ShortName: "DB", Birthstone: "April"},
        2: {Id: EMERALD, Name: "Emerald", ShortName: "EM", Birthstone: "May"},
        3: {Id: RUBY, Name: "Ruby", ShortName: "RY", Birthstone: "July"},
        4: {Id: SAPPHIRE_PINK, Name: "Pink Sapphire", ShortName: "SP", Birthstone: "September"},
        5: {Id: SAPPHIRE_YELLOW, Name: "Yellow Sapphire", ShortName: "SY", Birthstone: "September"},
        6: {Id: SAPPHIRE_BLUE, Name: "Blue Sapphire", ShortName: "SB", Birthstone: "September"}
    };
    var gemstoneClarityList = {
        0: {Id: GEM_CLARITY_ZERO, Name: "SI IJ", ShortName: "C0"},
        1: {Id: GEM_CLARITY_ONE, Name: "SI GH", ShortName: "C1"},
        2: {Id: GEM_CLARITY_TWO, Name: "VVS EF", ShortName: "C2"}
    };
    var gemstoneShapeList = {
        0: {Id: GEM_ROUND, Name: "Round"},
        1: {Id: GEM_SQUARE, Name: "Square"},
        2: {Id: GEM_MARQUISE, Name: "Marquise"},
        3: {Id: GEM_CUSHION, Name: "Cushion"},
        4: {Id: GEM_TRILLION, Name: "Trillion"},
        5: {Id: GEM_BAGUTTE, Name: "Bagutte"}
    };
    var gemstoneSettingList = {
        0: {Id: GEM_SET_PRONG, Name: "Prong"},
        1: {Id: GEM_SET_BEZEL, Name: "Bezel"},
        2: {Id: GEM_SET_PAVE, Name: "Pave"},
        3: {Id: GEM_SET_CHANNEL, Name: "Channel"}
    };
	
    var ALL_PRODUCTS = 0, DIAMOND_RING = 1, RING = 2, DIAMOND_PENDANT = 3, PENDANT = 4, DIAMOND_EARRING = 5, EARRING = 6, DIAMOND_BANGLE = 7, BANGLE = 8, DIAMOND_NOSEPIN = 9, NOSEPIN = 10, CHAIN = 11;
    var productCategoryList = {
        0: {Id: ALL_PRODUCTS, Name: "All products", ShortName: "AL"},
        1: {Id: DIAMOND_RING, Name: "Ring", ShortName: "DR"},
        2: {Id: RING, Name: "Ring", ShortName: "GR"},
        3: {Id: DIAMOND_PENDANT, Name: "Pendant", ShortName: "DP"},
        4: {Id: PENDANT, Name: "Pendant", ShortName: "GP"},
        5: {Id: DIAMOND_EARRING, Name: "Earring", ShortName: "DT"},
        6: {Id: EARRING, Name: "Earring", ShortName: "GT"},
        7: {Id: DIAMOND_BANGLE, Name: "Bangle", ShortName: "DB"},
        8: {Id: BANGLE, Name: "Bangle", ShortName: "GB"},
        9: {Id: DIAMOND_NOSEPIN, Name: "Nosepin", ShortName: "DN"},
        10: {Id: NOSEPIN, Name: "Nosepin", ShortName: "GN"},
        11: {Id: CHAIN, Name: "Chain", ShortName: "SC"}
    };

    this.resetCurVariant = function() {
        curVariant = null;
        curIndex = null;
    };
    this.getPrice = function(productInfo) {
        if(productInfo == null) {
            productInfo = curProductInfo;
        }
        curVariant = productVariants[ this.getIndex(productInfo) ];
        return curVariant.Price;
    };

    this.getIndex = function(productInfo) {
        if(productInfo == null) {
            productInfo = curProductInfo;
        }
        var tempIndex = '';
        var index = ''
                + productInfo.Id + '-'
                + metalColorList[ productInfo.PrimaryMetalColor ].ShortName;
        var clarityIndex = 0;
        if (productInfo.CenterStone != NO_STONE) {
            index += gemstoneList[ productInfo.CenterStone ].ShortName;
            //Only for diamond, other clarities are present
            clarityIndex = 0;
            if(productInfo.CenterStone == 0) {
                clarityIndex = productInfo.CenterStoneClarity;
            }
            tempIndex += gemstoneClarityList[ clarityIndex ].ShortName;
        }
        if (productInfo.AccentStoneOne != NO_STONE) {
            index += gemstoneList[ productInfo.AccentStoneOne ].ShortName;
            clarityIndex = 0;
            if(productInfo.AccentStoneOne == 0) {
                clarityIndex = productInfo.AccentStoneOneClarity;
            }
            tempIndex += gemstoneClarityList[ clarityIndex ].ShortName;
        }
        if (productInfo.AccentStoneTwo != NO_STONE) {
            index += gemstoneList[ productInfo.AccentStoneTwo ].ShortName;
            clarityIndex = 0;
            if(productInfo.AccentStoneTwo == 0) {
                clarityIndex = productInfo.AccentStoneTwoClarity;
            }
            tempIndex += gemstoneClarityList[ clarityIndex ].ShortName;
        }

        index += metalList[ productInfo.PrimaryMetal ].ShortName;
        if (tempIndex != '') {
            index += "-" + tempIndex;
        }
        if (productInfo == curProductInfo) {
            curIndex = index;
        }
        return index;
    };
    this.getRelativeImageUrl = function(view, productInfo) {
        if (productInfo == null) {
            productInfo = curProductInfo;
        }
        curVariant = productVariants[ this.getIndex(productInfo) ];
        var imageName = curVariant.ImageIndex + '-' + view + '.JPG';
        return productInfo.DirName + '/' + imageName;
    };
    this.getImage = function(view, productInfo) {
        if (productInfo == null) {
            productInfo = curProductInfo;
        }
        curVariant = productVariants[ this.getIndex(productInfo) ];
        var imageName = curVariant.ImageIndex + '-' + view + '.JPG';
        return configObj.DESIGNS_URL + '/'
                + productInfo.DirName + '/' + imageName;
    };
    this.getInspiringDesignMarkup = function(imageName, price, description) {
        var markup = '';
        markup += '<div class="inspiring_design">'
                +   '<img src="' + imageName + '" alt="" height="200">'
                +   '<div>' + description + '</div>'
                +   '<div class="inspiring_design_price">' 
                +        utilObj.formatMoney(price, true) 
                +   '</div>'
                + '</div>';
        return markup;
    };
    this.getProductInspirationMarkup = function() {
        var imageName, price, desc;
        var markup = '';
        var productInfo = jQuery.extend(true, {}, originalProductInfo);
        var prodIndex,variant;
        
        productInfo.PrimaryMetalColor = METAL_WHITE;
        productInfo.PrimaryMetal = GOLD_22KT;
        if (productInfo.CenterStone != NO_STONE) {
            productInfo.CenterStone = DIAMOND;
        }
        if (productInfo.AccentStoneOne != NO_STONE) {
            productInfo.AccentStoneOne = SAPPHIRE_BLUE;
        }
        prodIndex = this.getIndex(productInfo);
        variant = productVariants[ prodIndex ];
        imageName = this.getImage(1, productInfo);
        price = this.getPrice(productInfo);
        desc = this.getDescription(productInfo, variant)
        markup += this.getInspiringDesignMarkup(imageName, price, desc);

        productInfo.PrimaryMetalColor = METAL_YELLOW;
        productInfo.PrimaryMetal = GOLD_18KT;
        if (productInfo.CenterStone != NO_STONE) {
            productInfo.CenterStone = RUBY;
        }
        if (productInfo.AccentStoneOne != NO_STONE) {
            productInfo.AccentStoneOne = DIAMOND;
        }
        prodIndex = this.getIndex(productInfo);
        variant = productVariants[ prodIndex ];
        imageName = this.getImage(1, productInfo);
        price = this.getPrice(productInfo);
        desc = this.getDescription(productInfo, variant)
        markup += this.getInspiringDesignMarkup(imageName, price, desc);

        productInfo.PrimaryMetalColor = METAL_ROSE;
        productInfo.PrimaryMetal = GOLD_18KT;
        if (productInfo.CenterStone != NO_STONE) {
            productInfo.CenterStone = SAPPHIRE_BLUE;
        }
        if (productInfo.AccentStoneOne != NO_STONE) {
            productInfo.AccentStoneOne = SAPPHIRE_YELLOW;
        }
        prodIndex = this.getIndex(productInfo);
        variant = productVariants[ prodIndex ];
        imageName = this.getImage(1, productInfo);
        price = this.getPrice(productInfo);
        desc = this.getDescription(productInfo, variant)
        markup += this.getInspiringDesignMarkup(imageName, price, desc);

        productInfo.PrimaryMetalColor = METAL_WHITE;
        productInfo.PrimaryMetal = SILVER;
        if (productInfo.CenterStone != NO_STONE) {
            productInfo.CenterStone = EMERALD;
        }
        if (productInfo.AccentStoneOne != NO_STONE) {
            productInfo.AccentStoneOne = SAPPHIRE_PINK;
        }
        prodIndex = this.getIndex(productInfo);
        variant = productVariants[ prodIndex ];
        imageName = this.getImage(1, productInfo);
        price = this.getPrice(productInfo);
        desc = this.getDescription(productInfo, variant)
        markup += this.getInspiringDesignMarkup(imageName, price, desc);
        return markup;
    };

    this.getId = function() {
        return originalProductInfo.Id;
    };
    this.getEncId = function() {
        return originalProductInfo.EncId;
    };
    this.getName = function() {
        return originalProductInfo.Name;
    };
    this.getInfoForCart = function() {
        var productInfo = {};
        productInfo.Id = this.getId();
        productInfo.Name = this.getName();
        productInfo.Index = this.getIndex(null);
        productInfo.Price = this.getPrice();
        productInfo.Image = this.getRelativeImageUrl(1, null);
        productInfo.Desc = this.getDescription(null, null);
        productInfo.Qty = 1;
        return productInfo;
    };
    this.setCenterStone = function(stone) {
        curProductInfo.CenterStone = stone;
    };
    this.setAccentStoneOne = function(stone) {
        curProductInfo.AccentStoneOne = stone;
    };
    this.setAccentStoneTwo = function(stone) {
        curProductInfo.AccentStoneTwo = stone;
    };
    this.setPrimaryMetal = function(primaryMetal) {
        curProductInfo.PrimaryMetal = primaryMetal;
    };
    this.setPrimaryMetalColor = function(primaryMetalColor) {
        curProductInfo.PrimaryMetalColor = primaryMetalColor;
    };
    this.setCenterStoneClarity = function(clarity) {
        curProductInfo.CenterStoneClarity = clarity;
    };
    this.setAccentStoneOneClarity = function(clarity) {
        curProductInfo.AccentStoneOneClarity = clarity;
    };
    this.setAccentStoneTwoClarity = function(clarity) {
        curProductInfo.AccentStoneTwoClarity = clarity;
    };
    this.setCurrentView = function(view) {
        curView = view;
    };
    this.addThumbnailMarkup = function() {
        var numThumbnails = 5;
        var markup = '';
        if(curProductInfo.Category == DIAMOND_EARRING || curProductInfo.Category == DIAMOND_PENDANT) {
            numThumbnails = 3;
        }
        for(var i = 1; i <= numThumbnails; i++) {
            markup += '<img class="product_thumbnail" data-view="' + i
                 + '" src="" alt="" width="70" height="70">';
        }
        $('.product_thumbnails').html(markup);
    };
    this.showThumbnails = function() {
        $(".product_thumbnail").each(function() {
            $(this).attr('src', productObj.getImage( $(this).data('view'), null ) );
        });
    };
    this.requiresSizeSelector = function(category) {
        if (category == DIAMOND_RING || category == RING || 
                category == DIAMOND_BANGLE || category == BANGLE || 
                category == CHAIN) {
            return true;
        }
        return false;
    };
    this.showSizeSelector = function(category) {
        $('.size_selector').show();
    };
    this.showSharingOptions = function() {
        var tooltip = "tooltip";
        if(utilObj.isMobileBrowser()) {
            toolTip = "";
        }
        var sharingOptions = ''
                +'<div class="like_option like_option_design favorite clickable ' + tooltip + ' '
                + '" data-subjectid="' + originalProductInfo.EncId + '">'
                + 'Like'
                + '</div>'
                + '<div class="fb_share_design clickable" title="'
                +   'Share on Facebook" '
                + 'data-subjectid="' + originalProductInfo.EncId + '" data-picture="'
                +  'data.ImageName' + '">'
                + '<div>Share</div>'
                + '</div>';
        $('.product_sharing_options').html(sharingOptions);
    };
    this.showInitialProduct = function() {
        $('.delivery_month').text(originalProductInfo.DeliveryMonth);
        $('.delivery_day').text(originalProductInfo.DeliveryDay);
        if (this.requiresSizeSelector(originalProductInfo.Category)) {
            this.showSizeSelector(originalProductInfo.Category);
        }
        this.addThumbnailMarkup();
        this.showCurrent();
        this.showSharingOptions();
        clearTimeout(thumbnailTimer);
        productObj.showThumbnails();
        this.showCenterStoneChoices();
        this.hidePrimaryMetalChoices();
        this.updateSelectedChoices();
        this.bindInitialActions();
        if(!utilObj.isMobileBrowser()) {
            setTimeout(function() {
                $('.inspiring_designs').html(productObj.getProductInspirationMarkup());
            }, 5000);
        }
        if (originalProductInfo.AccentStoneOne != NO_STONE) {
            $('.accent_stone_one_wrapper').show();
            this.hideAccentStoneOneChoices();
        }
        if (originalProductInfo.AccentStoneTwo != NO_STONE) {
            $('.accent_stone_two_wrapper').show();
            this.hideAccentStoneTwoChoices();
        }
    };
    this.getDescription = function(product, variant) {
        if(product == null) {
            product = curProductInfo;
            variant = productVariants[ this.getIndex(product) ];
        }
        var desc = '';
        if (product.Category == DIAMOND_RING && variant.CenterStone == DIAMOND) {
            desc += gemstoneShapeList[product.CenterStoneShape].Name + " " + 
                    gemstoneList[variant.CenterStone].Name + " " + 
                    productCategoryList[product.Category].Name + " set in " + 
                    metalList[variant.PrimaryMetal].Karat + " " + 
                    metalColorList[variant.PrimaryMetalColor].Name + " " + 
                    metalList[variant.PrimaryMetal].Name
        } else if (product.Category == DIAMOND_RING || Product.Category == RING) {
            desc += gemstoneShapeList[product.CenterStoneShape].Name + " " +
                    gemstoneList[variant.CenterStone].Name + " " + 
                    metalList[variant.PrimaryMetal].Karat + " " + 
                    metalColorList[variant.PrimaryMetalColor].Name + " " + 
                    metalList[variant.PrimaryMetal].Name + " " + 
                    productCategoryList[product.Category].Name
        } else {
            desc += gemstoneShapeList[product.CenterStoneShape].Name + " " + 
                    gemstoneList[variant.CenterStone].Name + " " + 
                    metalList[variant.PrimaryMetal].Karat + " " + 
                    metalColorList[variant.PrimaryMetalColor].Name + " " + 
                    metalList[variant.PrimaryMetal].Name + " " + 
                    productCategoryList[product.Category].Name
        }

        if (variant.AccentStoneOne != NO_STONE) {
            desc += " with " + gemstoneList[variant.AccentStoneOne].Name
            if (variant.AccentStoneTwo != NO_STONE) {
                desc += " with " + gemstoneList[variant.AccentStoneTwo].Name
            }
        }
        return desc;
    };
    this.getDiscount = function() {
        return 10;
    };
    this.getOrigPrice = function(curPrice, discount) {
        return curPrice * 100 / (100 - discount);
    };

    this.showCurrent = function() {
        this.resetCurVariant();
        var curProdIndex = this.getIndex(curProductInfo);
        var curVariant = productVariants[ curProdIndex ];
        var curPrice = utilObj.formatMoney(curVariant.Price, true);
        var curDiscount = this.getDiscount();
        var origPrice = utilObj.formatMoney(this.getOrigPrice(curVariant.Price, curDiscount), false);
        var description = this.getDescription(curProductInfo, curVariant);
        $('.product_price').html(curPrice);
        $('.orig_product_price').html(origPrice);
        $('.product_discount').text(curDiscount);
        $('.product_desc').text(description);
        $('.product_code').text(curProdIndex);
        $('.product_name').text(curProductInfo.Name);
        var productImage = this.getImage(curView, curProductInfo);
        $('.product_image').attr('src', productImage);
        $('.zoomed_product').css('background-image', 'url("' + productImage + '")');
        clearTimeout(thumbnailTimer);
        thumbnailTimer = setTimeout(function() {
            productObj.showThumbnails();
        }, 10000);
        this.updateChosenNames();
    };
    this.updateChosenNames = function() {
        var stoneName = '';
        if(curProductInfo.CenterStone != NO_STONE) {
            stoneName = gemstoneList[curProductInfo.CenterStone].Name;
            if (curProductInfo.CenterStone == 0) {
                stoneName += " " + gemstoneClarityList[ curProductInfo.CenterStoneClarity ].Name;
            }
            $('.stone_' + CENTER_STONE + '_name').text(stoneName);
        }
        if(curProductInfo.AccentStoneOne != NO_STONE) {
            stoneName = gemstoneList[curProductInfo.AccentStoneOne].Name;
            if(curProductInfo.AccentStoneOne == 0) {
                stoneName += " " + gemstoneClarityList[ curProductInfo.AccentStoneOneClarity ].Name;
            }
            $('.stone_' + ACCENT_STONE_ONE + '_name').text(stoneName);
        }
        if (curProductInfo.AccentStoneTwo != NO_STONE) {
            stoneName = gemstoneList[curProductInfo.AccentStoneTwo].Name;
            if (curProductInfo.AccentStoneTwo == 0) {
                stoneName += " " + gemstoneClarityList[ curProductInfo.AccentStoneTwoClarity ].Name;
            }
            $('.stone_' + ACCENT_STONE_TWO + '_name').text(stoneName);
        }
        var metalName = metalList[curProductInfo.PrimaryMetal].Karat + " " + metalColorList[curProductInfo.PrimaryMetalColor].Name + " " + metalList[curProductInfo.PrimaryMetalColor].Name;
        $('.metal_primary_name').text(metalName);
    };
    this.bindInitialActions = function() {
        if (!utilObj.isMobileBrowser()) {
            $('.product_thumbnail').unbind("hover").hover(function() {
                var view = $(this).data('view');
                productObj.setCurrentView(view);
                productObj.switchView(view);
            });
            $('.product_image').unbind("hover").hover(function() {
                productObj.showZoomedProduct();
            }, function() {
                productObj.hideZoomedProduct();
            });
            $('.product_image').unbind("mousemove").mousemove(function(e) {
                productObj.updateZoomedLocation(e);
            });
        }
        $('.product_thumbnail').unbind("click").click(function() {
            var view = $(this).data('view');
            productObj.setCurrentView(view);
            productObj.switchView(view);
        });
        $('.center_stone_choices_toggle').unbind("click").click(function() {
            if ($('.center_stone_choices').is(':visible')) {
                productObj.hideCenterStoneChoices();
            } else {
                productObj.showCenterStoneChoices();
            }
        });
        $('.accent_stone_one_choices_toggle').unbind("click").click(function() {
            if ($('.accent_stone_one_choices').is(':visible')) {
                productObj.hideAccentStoneOneChoices();
            } else {
                productObj.showAccentStoneOneChoices();
            }
        });
        $('.accent_stone_two_choices_toggle').unbind("click").click(function() {
            if ($('.accent_stone_two_choices').is(':visible')) {
                productObj.hideAccentStoneTwoChoices();
            } else {
                productObj.showAccentStoneTwoChoices();
            }
        });
        $('.primary_metal_choices_toggle').unbind("click").click(function() {
            if ($('.primary_metal_choices').is(':visible')) {
                productObj.hidePrimaryMetalChoices();
            } else {
                productObj.showPrimaryMetalChoices();
            }
        });
    };
    this.bindActions = function() {
        if (!utilObj.isMobileBrowser()) {
            $('.customize_stone_option').unbind("hover").hover(function() {
                var stoneId = $(this).data('id');
                var stoneChoice = $(this).data('choice');
                if (stoneChoice == CENTER_STONE) {
                    productObj.setCenterStone(stoneId);
                } else if (stoneChoice == ACCENT_STONE_ONE) {
                    productObj.setAccentStoneOne(stoneId);
                } else if (stoneChoice == ACCENT_STONE_TWO) {
                    productObj.setAccentStoneTwo(stoneId);
                }
                productObj.showCurrent();
            }, function() {
                curProductInfo = jQuery.extend(true, {}, selectedProductInfo);
                productObj.showCurrent();
            });
            $('.customize_metal_option').unbind("hover").hover(function() {
                var metalId = $(this).data('metalid');
                var colorId = $(this).data('colorid');
                var metalName = $(this).data('name');
                $('.metal_primary_name').text(metalName);
                productObj.setPrimaryMetal(metalId);
                productObj.setPrimaryMetalColor(colorId);
                productObj.showCurrent();
            }, function() {
                curProductInfo = jQuery.extend(true, {}, selectedProductInfo);
                productObj.showCurrent();
            });
        }
        $('.customize_stone_option').unbind("click").click(function() {
            var stoneId = $(this).data('id');
            var stoneChoice = $(this).data('choice');
            if (stoneChoice == CENTER_STONE) {
                productObj.setCenterStone(stoneId);
            } else if (stoneChoice == ACCENT_STONE_ONE) {
                productObj.setAccentStoneOne(stoneId);
            } else if (stoneChoice == ACCENT_STONE_TWO) {
                productObj.setAccentStoneTwo(stoneId);
            }
            productObj.showCurrent();
            selectedProductInfo = jQuery.extend(true, {}, curProductInfo);
            productObj.updateSelectedChoices();
            if ($(this).hasClass('customize_diamond')) {
                productObj.showDiamondQualityChoices($(this).data('choice'));
            }
        });
        $('.customize_metal_option').unbind("click").click(function() {
            var metalId = $(this).data('metalid');
            var colorId = $(this).data('colorid');
            var metalName = $(this).data('name');
            $('.metal_primary_name').text(metalName);
            productObj.setPrimaryMetal(metalId);
            productObj.setPrimaryMetalColor(colorId);
            productObj.showCurrent();
            selectedProductInfo = jQuery.extend(true, {}, curProductInfo);
            productObj.updateSelectedChoices();
        });
    };
    this.showDiamondQualityChoices = function(stoneChoice) {
        var currentClarity = 0;
        if(stoneChoice == CENTER_STONE) {
            currentClarity = selectedProductInfo.CenterStoneClarity;
        } else if(stoneChoice == ACCENT_STONE_ONE) {
            currentClarity = selectedProductInfo.AccentStoneOneClarity;
        } else if(stoneChoice == ACCENT_STONE_TWO) {
            currentClarity = selectedProductInfo.AccentStoneTwoClarity;
        }
        var diamondChoiceMarkup = ''
            + '<div style="margin-bottom: 20px; letter-spacing: 0; color: #3E3E3E;">'
            +   'Choose the quality of your diamond'
            + '</div>'
            + '<div style="margin: auto; width: 280px;">'
            +   '<div class="diamond_quality_choice_wrapper">'
            +       '<div class="diamond_quality_choice customize_diamond" '
            +           'data-clarity="0" data-choice="' + stoneChoice + '">'
            +           '<div class="selected_stone '+ stoneChoice +'_0_selected"></div>'
            +       '</div>'
            +       '<div>SI IJ</div>'
            +   '</div>'
            +   '<div class="diamond_quality_choice_wrapper">'
            +       '<div class="diamond_quality_choice customize_diamond" '
            +           'data-clarity="1" data-choice="' + stoneChoice + '">'
            +           '<div class="selected_stone '+ stoneChoice +'_1_selected"></div>'
            +       '</div>'
            +       '<div>SI GH</div>'
            +   '</div>'
            +   '<div class="diamond_quality_choice_wrapper">'
            +       '<div class="diamond_quality_choice customize_diamond" '
            +           'data-clarity="2" data-choice="' + stoneChoice + '">'
            +           '<div class="selected_stone '+ stoneChoice +'_2_selected"></div>'       
            +       '</div>'
            +       '<div>VVS EF</div>'
            +   '</div>'
            + '</div>'
            + '<div class="like_link goudy_text clear" style="padding-top: 20px;">'
            + 'Not sure what\'s right for you? Let us help!'
            + '</div>';
        $('.diamond_selector_' + stoneChoice).html(diamondChoiceMarkup).show();
        $('.' + stoneChoice + '_' + currentClarity + '_selected').show();
        $('.diamond_quality_choice').unbind("click").click(function() {
            var stoneChoice = $(this).data('choice');
            productObj.selectDiamondQuality(stoneChoice, $(this).data('clarity'));
            $('.diamond_selector_' + stoneChoice).fadeOut();
            selectedProductInfo = $.extend(true, {}, curProductInfo);
            productObj.updateSelectedChoices();
            productObj.showCurrent();
        });
        if (!utilObj.isMobileBrowser()) {
            $('.diamond_quality_choice').unbind("hover").hover(function() {
                var stoneChoice = $(this).data('choice');
                productObj.selectDiamondQuality(stoneChoice, $(this).data('clarity'));
                productObj.showCurrent();
            }, function() {
                curProductInfo = jQuery.extend(true, {}, selectedProductInfo);
                productObj.showCurrent();
            });
        }
    };
    this.selectDiamondQuality = function(stoneChoice, stoneClarity) {
        if (stoneChoice == CENTER_STONE) {
            this.setCenterStoneClarity(stoneClarity);
        } else if (stoneChoice == ACCENT_STONE_ONE) {
            this.setAccentStoneOneClarity(stoneClarity);
        } else if (stoneChoice == ACCENT_STONE_TWO) {
            this.AccentStoneTwoClarity(stoneClarity);
        }
    };
    this.hideCenterStoneChoices = function() {
        $('.center_stone_choices_toggle').removeClass('down_triangle').addClass('left_triangle');
        $('.center_stone_choice_text').html(''
                + '<div class="selected_choice_desc">Primary stone</div>' 
                + '<span class="goudy_text selected_choice_text">' 
                + gemstoneList[selectedProductInfo.CenterStone].Name + '</span>');
        $('.center_stone_choices').hide();
    };
    this.hideAccentStoneOneChoices = function() {
        $('.accent_stone_one_choices_toggle').removeClass('down_triangle').addClass('left_triangle');
        $('.accent_stone_one_choice_text').html(''
                + '<div class="selected_choice_desc">Accent stone</div>' 
                + '<span class="goudy_text selected_choice_text">' 
                + gemstoneList[selectedProductInfo.AccentStoneOne].Name + '</span>');
        $('.accent_stone_one_choices').hide();
    };
    this.hideAccentStoneTwoChoices = function() {
        $('.accent_stone_two_choices_toggle').removeClass('down_triangle').addClass('left_triangle');
        $('.accent_stone_two_choice_text').html(''
                + '<div class="selected_choice_desc">Accent stone</div>' 
                + '<span class="goudy_text selected_choice_text">' 
                + gemstoneList[selectedProductInfo.AccentStoneTwo].Name + '</span>');
        $('.accent_stone_two_choices').hide();
    };
    this.hidePrimaryMetalChoices = function() {
        $('.primary_metal_choices_toggle').removeClass('down_triangle').addClass('left_triangle');
        var metalShortDesc = metalColorList[selectedProductInfo.PrimaryMetal].Name + " " + metalList[selectedProductInfo.PrimaryMetal].Name;
        $('.primary_metal_choice_text').html(''
                + '<div class="selected_choice_desc">Primary metal</div>' 
                + '<span class="goudy_text selected_choice_text">' 
                + metalShortDesc + '</span>');
        $('.primary_metal_choices').hide();
    };
    this.getStoneChoiceMarkup = function(choice) {
        var markup = ''
            + '<div class="stone_' + choice + '_0 customize_stone_option customize_diamond"'
            +     'data-id="0" data-choice="' + choice + '">'
            +   '<div class="selected_stone"></div>'
            + '</div>'
            + '<div class="stone_' + choice + '_4 customize_stone_option customize_pink_sapphire"'
            +     'data-id="4" data-choice="' + choice + '">'
            +   '<div class="selected_stone"></div>'
            + '</div>'
            + '<div class="stone_' + choice + '_5 customize_stone_option customize_yellow_sapphire"'
            +     'data-id="5" data-choice="' + choice + '">'
            +   '<div class="selected_stone"></div>'
            + '</div>'
            + '<div class="stone_' + choice + '_6 customize_stone_option customize_blue_sapphire"'
            +     'data-id="6" data-choice="' + choice + '">'
            +   '<div class="selected_stone"></div>'
            + '</div>'
            + '<div class="stone_' + choice + '_2 customize_stone_option customize_emerald"'
            +     'data-id="2" data-choice="' + choice + '">'
            +   '<div class="selected_stone"></div>'
            + '</div>'
            + '<div class="stone_' + choice + '_3 customize_stone_option customize_ruby"'
            +     'data-id="3" data-choice="' + choice + '">'
            +   '<div class="selected_stone"></div>'
            + '</div>'
            + '<div class="stone_' + choice + '_1 customize_stone_option customize_black_diamond"'
            +     'data-id="1" data-choice="' + choice + '">'
            +   '<div class="selected_stone"></div>'
            + '</div>'
            + '<div class="stone_' + choice + '_name stone_name cantarell_text">'
            + '</div>';
        return markup;
    };
    this.getMetalChoiceMarkup = function() {
        var markup = ''
            + '<div class="metal_primary_0_1 customize_metal_option customize_gold_22"'
            +     'data-metalid="0" data-colorid="1">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_1_1 customize_metal_option customize_gold_18"'
            +     'data-metalid="1" data-colorid="1">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_2_1 customize_metal_option customize_gold_14"'
            +     'data-metalid="2" data-colorid="1">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_0_0 customize_metal_option customize_white_gold_22"'
            +     'data-metalid="0" data-colorid="0">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_1_0 customize_metal_option customize_white_gold_18"'
            +     'data-metalid="1" data-colorid="0">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_2_0 customize_metal_option customize_white_gold_14"'
            +     'data-metalid="2" data-colorid="0">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_0_2 customize_metal_option customize_rose_22"'
            +     'data-metalid="0" data-name="22 Carat Rose Gold" data-colorid="2">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_1_2 customize_metal_option customize_rose_18"'
            +     'data-metalid="1" data-colorid="2">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_2_2 customize_metal_option customize_rose_14"'
            +     'data-metalid="2" data-colorid="2">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_3_0 customize_metal_option customize_silver"'
            +     'data-metalid="3" data-colorid="0">'
            +   '<div class="selected_metal"></div>'
            + '</div>'
            + '<div class="metal_primary_name metal_name cantarell_text">'
            + '</div>';
        return markup;
    };
    this.showCenterStoneChoices = function() {
        $('.center_stone_choice_text').text('Choose the primary stone');
        var stoneChoiceHtml = this.getStoneChoiceMarkup(CENTER_STONE);
        $('.center_stone_choices_toggle').removeClass('left_triangle').addClass('down_triangle');
        $('.center_stone_choices').html(stoneChoiceHtml).show();
        this.updateChosenNames();
        this.updateSelectedChoices();
        this.bindActions();
    };
    this.showAccentStoneOneChoices = function() {
        $('.accent_stone_one_choice_text').text('Choose the accent stone');
        var stoneChoiceHtml = this.getStoneChoiceMarkup(ACCENT_STONE_ONE);
        $('.accent_stone_one_choices_toggle').removeClass('left_triangle').addClass('down_triangle');
        $('.accent_stone_one_choices').html(stoneChoiceHtml).show();
        this.updateChosenNames();
        this.updateSelectedChoices();
        this.bindActions();
    };
    this.showAccentStoneTwoChoices = function() {
        $('.accent_stone_two_choice_text').text('Choose the accent stone');
        var stoneChoiceHtml = this.getStoneChoiceMarkup(ACCENT_STONE_TWO);
        $('.accent_stone_two_choices_toggle').removeClass('left_triangle').addClass('down_triangle');
        $('.accent_stone_two_choices').html(stoneChoiceHtml).show();
        this.updateChosenNames();
        this.updateSelectedChoices();
        this.bindActions();
    };
    this.showPrimaryMetalChoices = function() {
        $('.primary_metal_choice_text').text('Choose the primary metal');
        var metalChoiceHtml = this.getMetalChoiceMarkup();
        $('.primary_metal_choices_toggle').removeClass('left_triangle').addClass('down_triangle');
        $('.primary_metal_choices').html(metalChoiceHtml).show();
        this.updateChosenNames();
        this.updateSelectedChoices();
        this.bindActions();
    };
    //Relative positions: 
    //(0,0) -> (0,0)
    //(200,200) -> (-150, -150)
    //(400,400) -> (-300, -300)
    this.updateZoomedLocation = function(e) {
        var imageOffset = $('.product_image').offset();
        var imageX = e.pageX - imageOffset.left;
        var imageY = e.pageY - imageOffset.top;
        //var imageWidth = 400;
        //var zoomWidth = 600;
        //var realImageWidth = 900;
        var zoomPosX, zoomPosY;
        var multiplier = -0.75; //200 * -0.75 = -150
        zoomPosX = imageX * multiplier;
        zoomPosY = imageY * multiplier;
        
        $('.zoomed_product').css('background-position', zoomPosX + 'px ' + zoomPosY + 'px');
    };
    this.showZoomedProduct = function() {
        $('.zoomed_product').show();
    };
    this.hideZoomedProduct = function() {
        $('.zoomed_product').hide();
    };
    this.switchView = function(view) {
        var productImage = this.getImage(view, curProductInfo);
        $('.product_image').attr('src', productImage);
        $('.zoomed_product').css('background-image', 'url("' + productImage + '")');
    };
    this.updateSelectedChoices = function() {
        var centerStone = selectedProductInfo.CenterStone;
        var accentStoneOne = selectedProductInfo.AccentStoneOne;
        var accentStoneTwo = selectedProductInfo.AccentStoneTwo;
        var primaryMetal = selectedProductInfo.PrimaryMetal;
        var primaryMetalColor = selectedProductInfo.PrimaryMetalColor;
        $('.selected_stone').hide();
        $('.selected_metal').hide();
        if (centerStone != NO_STONE) {
            $('.stone_' + CENTER_STONE + '_' + centerStone + ' > .selected_stone').show();
        }
        if (accentStoneOne != NO_STONE) {
            $('.stone_' + ACCENT_STONE_ONE + '_' + accentStoneOne + ' > .selected_stone').show();
        }
        if (accentStoneTwo != NO_STONE) {
            $('.stone_' + ACCENT_STONE_TWO + '_' + accentStoneTwo + ' > .selected_stone').show();
        }
        $('.metal_primary_' + primaryMetal + '_' + primaryMetalColor + ' > .selected_metal').show();
    };
    
};

ProductList = function(list,choices) {
    var listData = JSON.parse(list);
    var products = new Array();
    var productUrlBase = configObj.BASE_URL + '/product';

    this.refinements = JSON.parse(choices);
    this.requestInProgress = false;
    this.allLotsDone = false;
    this.onScrollInProgress = false;
    this.lastListIndex = 0;
    this.numProducts = 0;
    
    var PRODUCT_LOT_SIZE = 9;
    for (var i = 0; i < listData.length; i++) {
        if(listData[i].ProductData.Id == 0) {
            continue;
        }
        this.numProducts++;
        if(i < PRODUCT_LOT_SIZE) {
            products[i] = new Product(listData[i]);
            this.lastListIndex = i;
        }
    }
    this.lastListIndex++;
    $('.num_products_found').text(this.numProducts);
    
    this.updateListData = function(curListData) {
        listData = curListData;
        this.lastListIndex = 0;
        this.numProducts = 0;
        for (var i = 0; i < listData.length; i++) {
            if (listData[i].ProductData.Id == 0) {
                continue;
            }
            this.numProducts++;
            if (i < PRODUCT_LOT_SIZE) {
                products[i] = new Product(listData[i]);
                this.lastListIndex = i;
            }
        }
        this.lastListIndex++;
        $('.num_products_found').text(this.numProducts);
    };
    this.createProducts = function(list) {
        var moreProducts = new Array();
        var productsCounter = 0;
        for (var i = 0; i < list.length && i < PRODUCT_LOT_SIZE; i++) {
            if (list[i].ProductData.Id != 0) {
                moreProducts[productsCounter] = new Product(list[i]);
                productsCounter++;
            }
        }
        return moreProducts;
    };
    this.showNextList = function() {
        if(productListObj.onScrollInProgress == true) {
            return;
        }
        productListObj.onScrollInProgress = true;
        var j = 0;
        var nextProductList = new Array();
        for (var i = this.lastListIndex; i < listData.length && j < PRODUCT_LOT_SIZE; i++,j++) {
            if (listData[i].ProductData.Id == 0) {
                continue;
            }
            nextProductList[j] = listData[i];
            this.lastListIndex = i;
        }
        this.lastListIndex++;
        var moreProducts = productListObj.createProducts(nextProductList);
        productListObj.show(moreProducts);
        if (this.lastListIndex == listData.length) {
            productListObj.allLotsDone = true;
        }
        productListObj.onScrollInProgress = false;
    };
    
    this.getMarkup = function(product) {
        var productUrl = productUrlBase + '/realShop?productCode=' + product.getIndex(null);
        var markup = ''
            + '<div class="product_wrapper" data-href="' + productUrl + '">'
            +   '<div style="width:240px;height:240px;">'
            +       '<img src="' + product.getImage(1, null) + '" width="240" height="240" alt="' 
            +       product.getName() + '">'
            +   '</div>'
            +   '<div class="like_link product_list_name">'
            +       product.getName()
            +   '</div>'
            +   '<div class="product_list_desc">' 
            +       '<span class="product_list_price">' 
            +           utilObj.formatMoney(product.getPrice(null), true)
            +       '</span>'
            +       ' ' + product.getDescription(null, null)
            +   '</div>'
            +   '<div class="product_wrapper_hover">'
            +       '<button class="light_pink_button customize_product" data-href="' + productUrl + '">'
            +       'Customize</button>'
            +   '</div>'
            + '</div>';
        return markup;
    };
    
    this.show = function(productsData) {
        if(productsData == null) {
            //For the initial load
            productsData = products;
            $(window).bind('scroll', this.onScroll);
        }
        if(productsData.length == 0) {
            productListObj.allLotsDone = true;
            return;
        }
        var markup = '';
        for(var i = 0; i < productsData.length; i++) {
            markup += this.getMarkup(productsData[i]);
        }
        $('.product_list').append(markup);
        utilObj.styleButtons();
        $('.product_wrapper').unbind("hover").hover( function() {
            $(this).children('.product_wrapper_hover').show();
        }, function() {
            $(this).children('.product_wrapper_hover').hide();
        });
        $('.product_wrapper').unbind('click').click( function(e) {
            e.preventDefault();
            window.location = $(this).data('href');
        });
        $('.customize_product').unbind('click').click( function(e) {
            e.preventDefault();
            window.location = $(this).data('href');
        });
        $('.refine_result').unbind('click').click(function() {
            productListObj.refineResult(this);
        });
    };
    
    this.afterNextList = function(result) {
        $('.loading').hide();
        if (result.Success == false) {
            productListObj.allLotsDone = true;
            return;
        }
        var curProductList = JSON.parse(result.Data.productList);
        productListObj.updateListData(curProductList);
        var curProducts = productListObj.createProducts(curProductList);
        productListObj.refinements = JSON.parse(result.Data.productRefinement);
        productListObj.show(curProducts);
        productListObj.requestInProgress = false;
    };
    this.getMarkupForSelected = function(refinementChoice) {
        var type2 = $(refinementChoice).data('type2');
        if(type2 == "undefined") {
            type2 = "";
        }
        var markup = ''
            + '<div class="selected_refinement ' + $(refinementChoice).data('type1') + '" data-type1="' 
            +   $(refinementChoice).data('type1') + '" data-type2="'
            +   type2 + '">'
            +   $(refinementChoice).data('name') 
            +   '<span class="deselect_refinement">X</span>'
            + '</div>';
        return markup;
    };
    this.refineResult = function(refinementChoice) {
        var type1 = $(refinementChoice).data('type1');
        var choice1 = $(refinementChoice).data('choice1');
        var type2 = $(refinementChoice).data('type2');
        var choice2 = $(refinementChoice).data('choice2');
        productListObj.refinements[type1] = choice1;
        if(type2 != "") {
          productListObj.refinements[type2] = choice2;
        }
        lastId = -1;
        var markup = productListObj.getMarkupForSelected(refinementChoice);
        $('.selected_refinement.' + type1).remove();
        $('.selected_refinements').append(markup);
        $('.deselect_refinement').unbind('click').click(function() {
            var type1 = $(this).parent().data('type1');
            var type2 = $(this).parent().data('type2');
            var refinementValue = -1;
            if(type1 == 'Category') {
                refinementValue = 0;
            }
            productListObj.refinements[type1] = refinementValue;
            if(type2 != "") {
                productListObj.refinements[type2] = -1;
            }
            $(this).parent().remove();
            $('.product_list').html('');
            lastId = -1;
            productListObj.fetchNextList();
        });
        $('.product_list').html('');
        $('.loading').show();
        productListObj.fetchNextList();
    };
 
    this.fetchNextList = function() {
        productListObj.lastListIndex = 0;
        productListObj.requestInProgress = true;
        productListObj.allLotsDone = false;
        $('.loading').show();
        var sortOrder = $('select.sorting_options').val();
        $.ajax({
            type: "POST",
            url: productUrlBase + "/getNextList",
            data: {
                'sort':     sortOrder,
                'category': productListObj.refinements.Category,
                'gemstone': productListObj.refinements.Gemstone,
                'setting':  productListObj.refinements.Setting,
                'priMetal': productListObj.refinements.PrimaryMetal,
                'priMetalColor': productListObj.refinements.PrimaryMetalColor,
                'priceMin': productListObj.refinements.PriceMin,
                'priceMax': productListObj.refinements.PriceMax
            },
            success: productListObj.afterNextList,
            dataType: 'json'
        });
    }; 
    this.onScroll = function() {
        if (productListObj.allLotsDone || productListObj.requestInProgress) {
            return;
        }
        // Check if we're within x pixels of the bottom edge of the browser window.
        var winHeight = window.innerHeight ? window.innerHeight : $(window).height(); // iphone fix
        var closeToBottom = ($(window).scrollTop() + winHeight > $(document).height() - 400);

        if (closeToBottom) {
            trackObj.count("action", configObj.CUR_PAGE, "page_scroll");
            productListObj.showNextList();
        }
    };
};

ProfanityFilter = function(externalSwears, externalValidWords) {
    var badWords,
            badWordRegexes = {},
            validWords,
            externalSwearsUrl = externalSwears,
            validWordsUrl = externalValidWords,
            alphaReplacements = {
                a: '(a|4|@)',
                b: '(b|8|3|ß|Β|β)',
                c: '(c|Ç|ç|¢|€|<)',
                d: '(d|Þ|þ|Ð|ð)',
                e: '(e|3|€|È|è|É|é|Ê|ê|∑)',
                f: '(f|ƒ)',
                g: '(g|6|9)',
                h: '(h)',
                i: '(i|1|∫|Ì|Í|Î|Ï|ì|í|î|ï)',
                j: '(j)',
                k: '(k|Κ|κ)',
                l: '(l|1|£|∫|Ì|Í|Î|Ï)',
                m: '(m)',
                n: '(n|η|Ν|Π)',
                o: '(o|0|Ο|ο|Φ|¤|°|ø)',
                p: '(p|ρ|Ρ|¶|þ)',
                q: '(q)',
                r: '(r|®)',
                s: '(s|5)',
                t: '(t|τ)',
                u: '(u|υ|µ)',
                v: '(v|υ|ν)',
                w: '(w|ω|ψ|Ψ)',
                x: '(x|Χ|χ)',
                y: '(y|¥|γ|ÿ|ý|Ÿ|Ý)',
                z: '(z)'
            };

    this.populateValidWords = function() {
        validWords = localStorage.getItem('validWords');
        if (validWords === null) {
            $.ajax({
                dataType: "json",
                url: validWordsUrl,
                data: {},
                success: this.storeValidWords
            });
        } else {
            validWords = JSON.parse(validWords);
        }
    };

    this.populateBadWords = function() {
        badWords = localStorage.getItem('badWords');
        if (badWords === null) {
            $.ajax({
                dataType: "json",
                url: externalSwearsUrl,
                data: {},
                success: this.storeBadWords
            });
        } else {
            badWords = JSON.parse(badWords);
        }
    };

    this.storeValidWords = function(result) {
        validWords = result;
        localStorage.setItem('validWords', JSON.stringify(validWords));
    };
    this.storeBadWords = function(result) {
        badWords = result;
        localStorage.setItem('badWords', JSON.stringify(badWords));
    };
    this.createBadWordRegexes = function() {
        for (var i = 0; i < badWords.length; i++) {
            var badWord = badWords[i];
            var badWordCombos = '';
            for (var j = 0; j < badWord.length; j++) {
                badWordCombos += alphaReplacements[badWord.charAt(j)];
            }
            badWordRegexes[badWord] = new RegExp(badWordCombos, 'i');
        }
    };
    this.filterBadWord = function(word) {
        if (validWords[word.toLowerCase()] == 1) {
            return word;
        }
        var filteredWord = word;
        if (badWordRegexes['shit'] == undefined) {
            this.createBadWordRegexes();
        }
        for (var i = 0; i < badWords.length; i += 1) {
            var re = badWordRegexes[ badWords[i] ];
            var rep = '*';
            if (re.test(filteredWord)) {
                filteredWord = filteredWord.replace(re, rep);
            }
        }
        return filteredWord;
    };
    this.filterText = function(text) {
        if (badWords === null) {
            this.populateBadWords();
            return text;
        }

        var words = text.split(" ");

        var filteredText = '';
        for (var wordIndex = 0; wordIndex < words.length; wordIndex++) {
            filteredText += this.filterBadWord(words[wordIndex]) + ' ';
        }

        return filteredText;
    };
};

RatingGame = function() {
    var reqForMysteryGift = 10;
    this.incrementRatingCount = function() {
        var ratingCount = localStorage.getItem("rated");
        if(ratingCount === null) {
            ratingCount = 1;
        } else {
            ratingCount = parseInt(ratingCount) + 1; 
        }
        
        localStorage.setItem("rated", ratingCount);
        this.updateState();
        
        if(ratingCount == reqForMysteryGift) {
            $('#rating_game_share').fadeIn();
			var cogsURL = configObj.BASE_URL + "/inspiration/playHot?cogs=quiz";
			facebookWrapperObj.publishCOGS('completed', 'hot_or_not', cogsURL);
        }
    };
    
    this.updateState = function() {
        var ratingCount = localStorage.getItem("rated");
        if(ratingCount === null) {
            $('#hot_progress_bar_text').text("Psst... There's a surprise in store for you!");
        } else {
            ratingCount = parseInt(ratingCount);
            var barWidth = (ratingCount * 100 / reqForMysteryGift);
            var barText = (reqForMysteryGift - ratingCount) + " more to go."
            if(ratingCount >= reqForMysteryGift) {
                barWidth = 100;
                barText = "";
                $('#hot_bar_wrapper').hide();
            }
            $('#hot_progress_bar').attr("style", "width: " + barWidth + "%;");
            $('#hot_progress_bar_text').text(barText);
        }
    };
};

Rate = function(top, current) {
    var topRated = null, currentPair = null;
    var hotOrNotBase = configObj.INSPIRATIONS_URL;
    if(top != "" && current != "") {
        topRated = JSON.parse(top);
        currentPair = JSON.parse(current);
    }
    
    this.showTopRated = function() {
        if(topRated == null) {
            return;
        }
        var topImagesHtml = '';
        for(var index=0; index < topRated.length; index++) {
            topImagesHtml += "<div class='top_rated_image hot_trends_link' data-src='top_rated'" 
                + " style='background-image: url("
                + hotOrNotBase + "/small/" + topRated[index].Filename + ");'>"
                + "</div>";
        }
        $('#top_rated_images').html(topImagesHtml);
    };
    
    this.showCurrentPair = function() {
        if(currentPair == null) {
            return;
        }
        $('#hot_img_left').html("<img class='hot_img'"
                + "src='" + hotOrNotBase + "/"
                + currentPair[0].Filename + "' alt='' height='400'>");
        $('#hot_left').attr('data-id', currentPair[0].EncId);
        $('#hot_img_right').html("<img class='hot_img'"
                + "src='" + hotOrNotBase + "/"
                + currentPair[1].Filename + "' alt='' height='400'>");
        $('#hot_right').attr('data-id', currentPair[1].EncId);
        ratingObj.enableRating();
        ratingObj.enableSkip();
    };
    
    this.enableRating = function() {
        $('.hot_img_wrapper').unbind('click').click(function() {
            trackObj.count("action", configObj.CUR_PAGE, "hotornot","rate");
            ratingObj.disableRating();
            ratingObj.rateCurrentPair(this);
        });
    };
    this.disableRating = function() {
        $('.hot_img_wrapper').unbind('click');
    };
    this.rateCurrentPair = function(winnerObj) {
        var winnerId = $(winnerObj).attr('data-id');
        var leftId = $('#hot_left').attr('data-id');
        var rightId = $('#hot_right').attr('data-id');
        var loserId;
        if(winnerId == leftId) {
            loserId = rightId;
        } else {
            loserId = leftId;
        }
        ratingGameObj.incrementRatingCount();
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/rate',
            data: {'winner': winnerId, 'loser': loserId},
            success: this.updateCurrentPair,
            dataType: 'json'
        });
        var storageKey = "lastVotedPhotoCOGS";
        var lastVoteCOGSAt = localStorage.getItem(storageKey);
        var curTimestamp = utilObj.getCurrentTimestamp();
        if (lastVoteCOGSAt === null || (curTimestamp - lastVoteCOGSAt) >= 86400) {
            localStorage.setItem(storageKey, curTimestamp);
            var cogsURL = configObj.BASE_URL + "/inspiration/playHot";
            facebookWrapperObj.publishCOGS('voted', 'photo', cogsURL);
        }
    };
    this.enableSkip = function() {
        $('#get_next_pair').unbind('click').click(function() {
            trackObj.count("action", configObj.CUR_PAGE, "hotornot","skip");
            ratingObj.disableSkip();
            ratingObj.skipCurrentPair();
        });
    };
    this.disableSkip = function() {
        $('#get_next_pair').unbind('click');
    };
    this.skipCurrentPair = function() {
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/rate/skip',
            data: {},
            success: this.updateCurrentPair,
            dataType: 'json'
        });
    };
    this.updateCurrentPair = function(result) {
        if(result.Success == false) {
            console.log(result);
            return;
        }
        currentPair = JSON.parse(result.Data.HotImages);
        ratingObj.showCurrentPair();
    };
};

Search = function(s, b) {
    var searchHost = s;
    var baseUrl = b;

    this.searchFor = function(queryString) {
        if (queryString == '') {
            return;
        }
        $.ajax({
            url: searchHost + 'select',
            data: {wt: 'json', q: queryString},
            success: this.searchResult,
            dataType: 'jsonp',
            jsonp: 'json.wrf'
        });
    };

    this.searchResult = function(data) {
        var results = 'No results found';
        if (data.response.numFound != 0) {
            results = '';
        }
        var docs = data.response.docs;
        for (var i = 0; i < docs.length; i++) {
            results += "<a href='" + baseUrl + "/jewelry/" + docs[i].name.replace(/ /g, '-') + "'>" + docs[i].name + '</a>'
                    + docs[i].description + '<br>';
        }
        $('#search_results').html(results);
    };
};

Sync = function() {
    var syncInterval = 15; //In seconds
    var fullFetchInterval = 600; //In seconds
    var maxSyncIndex = 3;
    //IsTemp: Whether or not the history of the actions needs to be maintained.
    //For eg: only the aggregate of views matters, not the individual views.
    var syncInfo = {
      0: {BaseKey: favoriteObj.getBaseKey(),  SyncUrl: favoriteObj.getSyncUrl(), SyncInProgress: false, IsTemp: favoriteObj.isTemp(), MergeForSync: favoriteObj.MergeForSync},
      1: {BaseKey: viewsObj.getBaseKey(),     SyncUrl: viewsObj.getSyncUrl(),    SyncInProgress: false, IsTemp: viewsObj.isTemp(),    MergeForSync: false},
      2: {BaseKey: cartObj.getBaseKey(),      SyncUrl: cartObj.getSyncUrl(),     SyncInProgress: false, IsTemp: cartObj.isTemp(),     MergeForSync: cartObj.MergeForSync},
      3: {BaseKey: expObj.getBaseKey(),       SyncUrl: expObj.getSyncUrl(),      SyncInProgress: true,  IsTemp: expObj.isTemp(),      MergeForSync: false}
    };

    var backendSync = window.setInterval(function() {
            syncObj.syncToBackend();
            syncObj.fetchFromBackend();
        }, syncInterval * 1000);

    this.isSyncInProgress = function(index) {
        return syncInfo[index].SyncInProgress;
    };

    this.isSyncEnabled = function() {
        if(!userObj.isLoggedIn()) {
            return false;
        }
        return true;
    };
    this.isFetchEnabled = function() {
        if(!userObj.isLoggedIn()) {
            return false;
        }
        //if(userObj.isTester()) {
        //So that experiment data is not overwritten continuosly
        //return false;
        //}
        return true;
    };
    
    this.syncToBackend = function() {
        if( this.isSyncEnabled() === false ) {
            this.stopSync();
            return;
        }
        for (var i = 0; i <= maxSyncIndex; i++) {
            if (syncInfo[i].SyncInProgress) {
                continue;
            }
            var localStorageKey = utilObj.getLocalKey(syncInfo[i].BaseKey);
            var localData = localStorage.getItem(localStorageKey);
            var tempStorageKey = utilObj.getTempKey(syncInfo[i].BaseKey);
            var tempData = localStorage.getItem(tempStorageKey);
            var dataToSync = localData;
            if (dataToSync === null) {
                dataToSync = tempData;
            }
            if (tempData !== null && localData !== null) {
                dataToSync = syncInfo[i].MergeForSync(dataToSync,tempData);
            }
            if(dataToSync === null) {
                dataToSync = '';
            }

            var removeStorageKey = utilObj.getRemoveKey(syncInfo[i].BaseKey);
            var dataToRemove = localStorage.getItem(removeStorageKey);
            if (dataToRemove === null) {
                dataToRemove = '';
            }
            if (dataToSync != '' || dataToRemove != '') {
                if (tempData !== null) {
                    //Move temp data to local data
                    localStorage.setItem(localStorageKey, dataToSync);
                    localStorage.removeItem(tempStorageKey);
                }
                //Sync the local data to backend
                syncInfo[i].SyncInProgress = true;
                $.ajax({
                    type: "POST",
                    url:  configObj.BASE_URL + '/sync' + syncInfo[i].SyncUrl,
                    data: {'index': i, 'localData': dataToSync, 'removedData': dataToRemove},
                    success: this.syncDone,
                    dataType: 'json'
                });
            }
        }
    };

    this.syncDone = function(result) {
        var index = result.Data.index;
        if (result.Success === false) {
            console.log(result);
            syncObj.stopSync(); //Something went wrong. Stop syncs for this session.
            return;
        }

        var storageKey = utilObj.getLocalKey(syncInfo[index].BaseKey);
        var localData = localStorage.getItem(storageKey);
        if (localData === null) {
            console.log("Why is localData null after sync success??");
        }
        var tempData = localStorage.getItem(utilObj.getTempKey(syncInfo[index].BaseKey));
        if (tempData === null) {
          localStorage.removeItem(storageKey);
        } else {
          localStorage.setItem(storageKey, tempData);
        }
        localStorage.removeItem(utilObj.getTempKey(syncInfo[index].BaseKey));
        if (syncInfo[index].IsTemp === false) {
            localStorage.setItem(syncInfo[index].BaseKey, result.Data.syncedData);
            localStorage.removeItem(utilObj.getRemoveKey(syncInfo[index].BaseKey));
        }
        syncInfo[index].SyncInProgress = false;
    };

    this.fetchFromBackend = function() {
        if(!this.isFetchEnabled()) {
            return;
        }
        //At any point, the base key has the data which is already present in the backend.
        //So overwriting the base key by fetching from backend occasionally is fine.
        var lastFetchAt = localStorage.getItem('lastFetch');
        var currentTimestamp = utilObj.getCurrentTimestamp();
        if (lastFetchAt !== null
                && (currentTimestamp - lastFetchAt) < fullFetchInterval) {
            return;
        }
        
        for (var i = 0; i <= maxSyncIndex; i++) {
            if (syncInfo[i].IsTemp === true || syncInfo[i].SyncInProgress) {
                continue;
            }
            syncInfo[i].SyncInProgress = true;
            $.ajax({
                type: "POST",
                url: configObj.BASE_URL + '/fetch' + syncInfo[i].SyncUrl,
                data: {'index': i},
                success: this.fetchDone,
                dataType: 'json'
            });
        }
    };

    this.fetchDone = function(result) {
        var index = result.Data.index;
        if (result.Success == false) {
            console.log(result);
            syncObj.stopSync();
            return;
        }
        localStorage.setItem(syncInfo[index].BaseKey, result.Data.syncedData);
        var currentTimestamp = utilObj.getCurrentTimestamp();
        localStorage.setItem('lastFetch', currentTimestamp);
        syncInfo[index].SyncInProgress = false;
    };

    this.stopSync = function() {
        clearInterval(backendSync);
    };
};

Tag = function(allTags) {
    var tagList = null;
    var categoryList = new Array();
    if(allTags != undefined && allTags != "") {
        tagList = JSON.parse(allTags);
    }
    
    this.getAllTags = function() {
        return tagList;
    };
    
    this.populateCategoryLists = function() {
        if(categoryList['metal'] != undefined) {
            return;
        }
        for(var i = 0; i < tagList.length; i++) {
            var category = tagList[i].Category;
            if(category == "") {
                continue;
            }
            if(categoryList[category] == undefined) {
                categoryList[category] = new Array();
            }
            categoryList[category].push(tagList[i]);
        }
    };
    
    this.getCategoryList = function() {
        this.populateCategoryLists();
        return categoryList;
    };
    this.showCategoryList = function(whereToShow) {
        this.populateCategoryLists();
        var categoriesInfo = '';
        for (var category in categoryList) {
            categoriesInfo += "<div>" + category + ": ";
            for(var i=0; i < categoryList[category].length; i++) {
                categoriesInfo += categoryList[category][i].Tag + ",";
            }
            categoriesInfo += "</div>";
        }
        $('#' + whereToShow).html(categoriesInfo);
    };
    
    this.getAllTagsForCategory = function(category) {
        if(categoryList[category] == undefined) {
            this.populateCategoryLists();
        }
        return categoryList[category];
    };
};

Track = function() {
    var trackBaseUrl = configObj.BASE_URL + '/track';
    this.count = function(kingdom, phylum, cla, family, genus, species, incrementBy, order) {
        var userId = userObj.getUserTrackingId();
        if (kingdom == undefined || phylum == undefined) {
            return;
        }
        if (cla == undefined) {
            cla = "";
        }
        if (family == undefined) {
            family = "";
        }
        if (genus == undefined) {
            genus = "";
        }
        if (species == undefined) {
            species = "";
        }
	if (incrementBy == undefined) {
            incrementBy = 1;
	}
        if (order == undefined) {
            order = "";
        }

        var os = utilObj.getOSInfo();
        var browser = utilObj.getBrowserInfo();
        var tester = utilObj.isTester();
        var params = {
            "userId":   userId,
            "kingdom":  kingdom,
            "phylum":   phylum,
            "class":    cla,
            "family":   family,
            "genus":    genus,
            "species":  species,
            "order":    order,
            "os":       os,
            "browser":  browser,
            "tester":   tester,
            "incrementBy": incrementBy,
            "width":    window.innerWidth,
            "height":   window.innerHeight
        };
        $.ajax({
            type: "POST",
            url: trackBaseUrl,
            data: params
        });
    };
};

User = function(loggedIn,data,ref) {
    var loginStatus = loggedIn;
    var userData = JSON.parse(data);
    this.referrer = ref;
    
    if(userData.Email != "") {
        localStorage.setItem('userData', JSON.stringify(userData));
    }
    this.askEmail = function() {
        if (utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/askEmail?jref=' + userObj.referrer + '&utm_source=' + configObj.CUR_PAGE;
        } else {
            utilObj.createPopup(configObj.BASE_URL + '/user/askEmail',
                {}, 500, 350, false);
        }
    };
    this.isLoggedIn = function() {
        return loginStatus;
    };
    
    this.storeUserData = function(data) {
        localStorage.setItem('userData', JSON.stringify(data));
    };
    
    this.getUserData = function() {
        if(userData.Email == "") {
            var localUserData = localStorage.getItem('userData');
            if(localUserData !== null) {
                userData = JSON.parse(localUserData);
            }
        }
        return userData;
    };
    
    this.getUserId = function() {
        userData = this.getUserData();
        if(userData.EncId == "") {
            return 0;
        }
        return userData.EncId;
    };
    this.getUserTrackingId = function() {
        var id = this.getUserId();
        var tempId = localStorage.getItem('tempId');
        if (tempId !== null && id != 0) {
            localStorage.removeItem('tempId');
            trackObj.count("debug", "temp_id", tempId);
        }
        if(id != 0) {
            return id;
        }
        
        if (tempId === null) {
            var tempId = "rand_" + utilObj.getRandomInt(1, 1000000);
            localStorage.setItem('tempId', tempId);
        }
        return tempId;
    };
    
    this.getReferralCount = function() {
        userData = this.getUserData();
        if(userData.ReferralCount == "") {
            return 0;
        }
        return userData.ReferralCount;
    };
    
    this.getFbid = function() {
        userData = this.getUserData();
        return userData.Fbid;
    };
    
    this.getEmail = function() {
        userData = this.getUserData();
        return userData.Email;
    };
    this.gmailImportStatus = function() {
        userData = this.getUserData();
        return userData.GmailImport;
    };
    
    this.yahooImportStatus = function() {
        userData = this.getUserData();
        return userData.YahooImport;
    };
    
    this.isTester = function() {
        var tester = localStorage.getItem('tester');
        if (tester == 1) {
            return true;
        }
        return false;
    };
    
    this.submitEmail = function(email) {
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/user/create',
            data: {'login_email': email, 'ref': userObj.referrer},
            success: this.afterSubmitEmail,
            dataType: 'json'
        });
    };
    this.afterSubmitEmail = function(result) {
        if(result.Success == true) {
            if(utilObj.isMobileBrowser()) {
                window.location = configObj.BASE_URL + '/?jref=' + userObj.referrer + '&utm_source=' + configObj.CUR_PAGE;
            } else {
                utilObj.reloadPage();
            }
        } else {
            alert("Unable to create account." + result.Data.message);
        }
    };
    this.completeSignup = function(name, password, verifyPassword) {
        if(!utilObj.isMobileBrowser()) {
            utilObj.closePopup();
        }
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/user/completeSignup',
            data: {'username': name, 'password': password, 'verifyPassword': verifyPassword},
            success: this.showInitialCongrats,
            dataType: 'json'
        });
    };
    this.signupCompletionReqd = function() {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/signupCompletionReqd?utm_source=' + configObj.CUR_PAGE;
        } else {
            trackObj.count("view", configObj.CUR_PAGE, "enter_details");
            utilObj.createPopup(configObj.BASE_URL + '/user/signupCompletionReqd',
                {}, 500, 370, false);
        }
    };
    this.showChangePassword = function() {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/showChangePassword?utm_source=' + configObj.CUR_PAGE;
        } else {
            trackObj.count("view", configObj.CUR_PAGE, "change_pwd");
            utilObj.createPopup(configObj.BASE_URL + '/user/showChangePassword',
                {}, 500, 300, true);
        }
    };
    this.showInitialCongrats = function() {
        trackObj.count("view", configObj.CUR_PAGE, "credit_500");
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/showInitialCongrats?utm_source=' + configObj.CUR_PAGE;
        } else {
            utilObj.createPopup(configObj.BASE_URL + '/user/showInitialCongrats',
                {}, 500, 280, true);
			var cogsURL = configObj.BASE_URL + "?cogs=cash";
			facebookWrapperObj.publishCOGS('received', 'five_hundred_cash', cogsURL);
        }
    };
    this.showFBConnect = function() {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/askFBConnect?jref=' + userObj.referrer + '&utm_source=' + configObj.CUR_PAGE;
        } else {
            trackObj.count("view", configObj.CUR_PAGE, "ask_fb_connect");
            utilObj.createPopup(configObj.BASE_URL + '/user/askFBConnect',
                {}, 500, 300, true);
        }
    };
    this.showForgotPasswordPopup = function() {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/forgotPassword?jref=' + userObj.referrer + '&utm_source=' + configObj.CUR_PAGE;
        } else {
			utilObj.createPopup(configObj.BASE_URL + '/user/forgotPassword',
                {}, 400, 250, true);
		}
	};
    this.showSignupPopup = function() {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/signup?jref=' + userObj.referrer + '&utm_source=' + configObj.CUR_PAGE;
        } else {
            utilObj.createPopup(configObj.BASE_URL + '/user/signup',
                {}, 400, 350, true);
        }
    };
    
    this.showLoginPopup = function() {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/login?jref=' + userObj.referrer + '&utm_source=' + configObj.CUR_PAGE;
        } else {
            utilObj.createPopup(configObj.BASE_URL + '/user/login',
                {}, 400, 430, true);
        }
    };
    
    this.showInviteFriends = function(network) {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/showInviteFriends?network=' + network 
                    + '&utm_source=' + configObj.CUR_PAGE;
        } else {
            trackObj.count("view", configObj.CUR_PAGE, "invite", network);
            utilObj.createPopup(configObj.BASE_URL + '/user/showInviteFriends',
                {'network': network}, 800, 380, true);
        }
    };
    
    this.storeFBRequest = function(reqId, recepients) {
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/user/storeFBRequest',
            data: {'reqId': reqId, 'recepients': recepients}
        });
    };
    
    this.invitesSent = function(count) {
        if(utilObj.isMobileBrowser()) {
            window.location = configObj.BASE_URL + '/user/invitesSent?count=' + count + '&utm_source=' + configObj.CUR_PAGE;
        } else {
            trackObj.count("view", configObj.CUR_PAGE, "invitesSent");
            utilObj.createPopup(configObj.BASE_URL + '/user/invitesSent',
                {"count": count}, 800, 380, true);
        }
    };
    
    this.populateContacts = function(emailDataStr) {
        var htmlStr = "";
        var emailData = JSON.parse(emailDataStr);
        var emailList = emailData.Emails.split(";");
        var namesList = emailData.Names.split(";");
        var totalContacts = 0;
        var contactsLimit = 100;
        if(!utilObj.isMobileBrowser()) {
            contactsLimit = 1000;
        }
        for( var i = 0; i < emailList.length && i < contactsLimit; i++) {
            if(emailList[i] == "") {
                continue;
            }
            //namesList[i] = namesList[i].replace(/"/g, '\"');
            totalContacts++;
            htmlStr += ''
                + '<div class="contact_info">'
                +   '<div class="contacts_selected_checkbox">'
                +        '<input type="checkbox" value="' + emailList[i] 
                +           '" name="contact">'
                +   '</div>'
                +   '<div class="contacts_selector_info">'
                +   	'<div class="contacts_selector_name" title="' + namesList[i] + '">'
                +       	namesList[i]
                +   	'</div>'
                +   	'<div class="contacts_selector_email" title="' + emailList[i] + '">'
                +       	emailList[i]
                +   	'</div>'
                +   '</div>'
                + '</div>';
        }
        $('.contacts_list').html(htmlStr);
        $('#total_contacts').text(totalContacts);
    };
    this.removeAddress = function(addressId) {
        if(this.isLoggedIn == false) {
          return;
        }
        $.ajax({
            type: "POST",
            url: configObj.BASE_URL + '/user/removeAddress',
            data: {'addressId': addressId},
            dataType: 'json'
        });
    };
};

Util = function() {
    var closeButtonVisible;
    var popupWidth;
    var popupHeight;
   
    this.bindLinks = function() {
        $('.play_hot_link').unbind('click').click(function(e) {
            e.preventDefault();
            var utm_source = $(this).data('src');
            var utm_medium = $(this).data('medium');
            if (utm_medium == undefined) {
                utm_medium = "";
            }
            window.location = configObj.BASE_URL + '/inspiration/playHot?jref=' + userObj.referrer + '&utm_source=' + utm_source + '&utm_medium=' + utm_medium;
        });
        $('.win_pendant_link').unbind('click').click(function(e) {
            e.preventDefault();
            var utm_source = $(this).data('src');
            var utm_medium = $(this).data('medium');
            if (utm_medium == undefined) {
                utm_medium = "";
            }
            window.location = configObj.BASE_URL + '/contest/winPendant?jref=' 
                    + userObj.referrer + '&utm_source=' + utm_source + '&utm_medium=' 
                    + utm_medium;
        });
        $('.earn_cash_link').unbind('click').click(function(e) {
            e.preventDefault();
            var utm_source = $(this).data('src');
            var utm_medium = $(this).data('medium');
            if (utm_medium == undefined) {
                utm_medium = "";
            }
            window.location = configObj.BASE_URL + '/user/earnCash?jref=' 
                    + userObj.referrer + '&utm_source=' + utm_source + '&utm_medium=' 
                    + utm_medium;
        });
        $('.hot_trends_link').unbind('click').click(function(e) {
            e.preventDefault();
            var utm_source = $(this).data('src');
            var utm_medium = $(this).data('medium');
            if (utm_medium == undefined) {
                utm_medium = "";
            }
            window.location = configObj.BASE_URL + '/inspiration/hotTrends?jref=' 
                    + userObj.referrer + '&utm_source=' + utm_source + '&utm_medium=' 
                    + utm_medium;
        });
        $('.shop_link').unbind('click').click(function(e) {
            e.preventDefault();
            var utm_source = $(this).data('src');
            var utm_medium = $(this).data('medium');
            if(utm_medium == undefined) {
                utm_medium = "";
            }
            window.location = configObj.BASE_URL + '/product/shop?jref=' 
                    + userObj.referrer + '&utm_source=' + utm_source + '&utm_medium=' 
                    + utm_medium;
        });
        $('.cart').unbind('click').click(function(e) {
            e.preventDefault();
            window.location = configObj.BASE_URL + '/user/viewCart';
        });
        $('.checkout').unbind('click').click(function(e) {
            e.preventDefault();
            window.location = configObj.BASE_URL + '/user/checkout';
        });
    };
    this.createPopup = function(urlToHit, params, width, height, showClose) {
        this.setPopupParams(showClose, width, height);
        $.ajax({
            type: "POST",
            url: urlToHit,
            data: params,
            success: this.showPopup,
            dataType: 'html',
            error: function(w, t, f) {
                $('#popup_overlay').hide();
                //alert(w + "\n" + t + "\n" + f);
                console.log(w,t,f);
            }
        });
    };
    this.setPopupParams = function(closeVisible, width, height) {
        $('#popup_overlay').unbind('click').show();
        if(closeVisible) {
            $('#popup_overlay').click(function() {
               utilObj.closePopup(); 
            });
        }
        $('#popup').addClass('bounce_show');
        closeButtonVisible = closeVisible;
        popupWidth = width;
        popupHeight = height;
    };
    this.closePopup = function() {
        $('#popup').addClass('bounce_hide').removeClass('bounce_show').html('');
        $('#popup').attr('style', 'display: block; width: 0; height: 0;').fadeOut(500, function() {
            $(this).removeClass('bounce_hide').attr('style', '');
        });
        $('#popup_overlay').fadeOut();
    };
    this.showPopup = function(data) {
        var initialStyle = 'display: block; width: ' + popupWidth + 'px; height: ' + popupHeight +
                'px;margin-left: -' + (popupWidth / 2.0) + 'px; margin-top: -' +
                (popupHeight / 2.0) + 'px;'
        if (closeButtonVisible == true) {
            data = '<div id="popup_close"></div>' + data;
            $('#popup').show(10, function() {
                $(this).attr('style', initialStyle).html(data);
                $('#popup_close').show().click(function() {
                    utilObj.closePopup();
                });
            });
        } else {
            $('#popup').show(10, function() {
                $(this).attr('style', initialStyle).html(data);
            });
        }
    };
    this.showErrorMessage = function(msg) {
        var popupWidth = 500;
        var popupHeight = 100;
        var initialStyle = 'display: block; width: ' + popupWidth + 'px; height: ' + popupHeight +
                'px;margin-left: -' + (popupWidth / 2.0) + 'px; margin-top: -' +
                (popupHeight / 2.0) + 'px;'
        msg = '<div id="alert_close"></div>' + msg;
        $('#alert_box').show(10, function() {
            $(this).attr('style', initialStyle).html(msg);
            $('#alert_close').show().click(function() {
                $('#alert_box').hide();
            });
        });  
    };
    this.reloadPage = function() {
        var curLocation = window.location.pathname;
        if (curLocation.indexOf("logout") == -1) {
            window.location.reload();
        } else {
            window.location = configObj.BASE_URL;
        }
    };
    this.getFirstName = function(name) {
        var splitName = name.split(" ");
        return splitName[0];
    };
    this.styleTooltips = function(element) {
        element = typeof element !== 'undefined' ? element : '.tooltip';
        $(element).tooltip({
            position: {
                my: "center bottom-20",
                at: "center top",
                using: function(position, feedback) {
                    $(this).css(position);
                    $("<div>")
                            .addClass("arrow")
                            .addClass(feedback.vertical)
                            .addClass(feedback.horizontal)
                            .appendTo(this);
                }
            }
        });
    };
    
    this.styleInputsAndButtons = function() {
        this.styleTextAreas();
        this.styleButtons();
        this.styleInputAreas();
        this.styleTooltips();
    };

    this.styleTextAreas = function() {
        $('input[type=textarea]').button().css({
            'font': 'inherit',
            'color': 'inherit',
            'text-align': 'left',
            'outline': 'none',
            'cursor': 'text',
            'border': '1px solid rgb(200, 200, 200)',
            'margin-left': '0',
            'letter-spacing': '1px',
            'line-height': '14px',
            'margin-top': '5px',
            'font-size': '13px',
            'padding':  '0.4em',
            'font-family': '"Lucida Grande","Lucida Sans Unicode",Helvetica,Arial,Verdana,sans-serif'
        });
    };

    this.styleButtons = function() {
        $("input[type=submit], button").button().click(function(event) {
            event.preventDefault();
        }).css({
            'color': 'white',
            'font-weight': 'normal',
            'border-radius': '0px',
            'margin': '0',
            'letter-spacing': '1px',
            'line-height': '22px',
            'height': '37px',
            'font-size': '16px',
            'font-family': '"Lucida Grande","Lucida Sans Unicode",Helvetica,Arial,Verdana,sans-serif'
        });
        $(".style_as_button").button().css({
            'color': 'white',
            'font-weight': 'normal',
            'border-radius': '0px',
            'letter-spacing': '2px',
            'line-height': '22px',
            'height': '34px',
            'font-size': '16px',
            'font-family': '"Lucida Grande","Lucida Sans Unicode",Helvetica,Arial,Verdana,sans-serif'
        });
    };

    this.styleInputAreas = function() {
        $('input:text, input:password, input[type=email]').button().css({
            'text-align': 'left',
            'outline': 'none',
            'cursor': 'text',
            'border': '1px solid rgb(200, 200, 200)',
            'margin-left': '0',
            'padding':  '0.4em',
            'height': '37px',
            'font-family': '"Cantarell", sans-serif'
        });
    };

    this.getCurrentTimestamp = function() {
        return Math.round(new Date().getTime() / 1000);
    };

    this.convertTimestampToText = function(timestamp) {
        var currentTimestamp = this.getCurrentTimestamp();
        var diff = currentTimestamp - timestamp;
        if (diff <= 0) {
            return 'just now';
        } else if (diff < 60) {
            return 'few seconds ago';
        } else if (diff < 1200) {
            return 'few minutes ago';
        } else if (diff < 3600) {
            return 'about an hour ago';
        } else if (diff < 86400) {
            return 'few hours ago';
        } else if (diff < 2592000) {
            return 'few days ago';
        } else {
            return 'few months back';
        }
    };

    this.enablePreviews = function() {
        xOffset = 0;
        yOffset = -100;

        $("div.preview_enabled").mouseenter(function(e) {
            var imgUrl = $(this).data('imgurl');
            $(this).append("<img id='preview' src='" + imgUrl
                    + "' alt='Image preview' />");
            $("#preview")
                    .css("top", (yOffset) + "px")
                    .css("left", (xOffset) + "px")
                    .fadeIn("fast");
        });
        $("div.preview_enabled").mouseleave(function() {
            $("#preview").remove();
        });
    };

    this.getLocalKey = function(baseKey) {
        return baseKey + 'Local';
    };

    this.getTempKey = function(baseKey) {
        return baseKey + 'Temp';
    };

    this.getRemoveKey = function(baseKey) {
        return baseKey + 'Remove';
    };
    
    this.getCacheKey = function(key) {
        return 'jools_' + key;
    };
    
    this.getRandomInt = function(min, max) {
        return Math.floor((Math.random() * max) + min);
    };

    this.getBrowserInfo = function() {
        var userAgent = navigator.userAgent;
        var browser = localStorage.getItem('browser');
        if(browser !== null) {
            return browser;
        }
        var chrome = "Chrome";
        var firefox = "Firefox";
        var safari = "Safari";
        var opera = "Opera";
        var dolfin = "Dolfin";
        var other = "Other";
        
        if (userAgent.indexOf(chrome) !== -1) {
            browser = chrome;
        } else if (userAgent.indexOf(firefox) !== -1) {
            browser = firefox;
        } else if (userAgent.indexOf(safari) !== - 1) {
            browser = safari;
        } else if (userAgent.indexOf(opera) !== - 1) {
            browser = opera;
        } else if (userAgent.indexOf(dolfin) !== - 1) {
            browser = dolfin;
        } else {
            browser = other;
        }
        //localStorage.setItem('browser', browser);
        return browser;
    };
    
    this.isMobileBrowser = function() {
        //return true;
        var userAgent = navigator.userAgent;
        if (userAgent.indexOf("Mobile") !== -1) {
            return true;
        } else if (userAgent.indexOf("Phone") !== -1) {
            return true;
        } else if (userAgent.indexOf("Android") !== -1) {
            return true;
        }
        return false;
    };
    
    this.getOSInfo = function() {
        var userAgent = navigator.userAgent;
        var os = localStorage.getItem('os');
        if (os !== null) {
            return os;
        }
        
        var linux = "Linux";
        var ipad = "iPad";
        var iphone = "iPhone";
        var android = "Android";
        var windowsPhone = "Windows Phone";
        var windows = "Windows";
        var mac = "Macintosh";
        var blackberry = "BlackBerry";
		var bada = "Bada";
        var other = "Other";
        
        if (userAgent.indexOf(ipad) != - 1) {
            os = ipad;
        } else if (userAgent.indexOf(iphone) != - 1) {
            os = iphone;
        } else if (userAgent.indexOf(android) != - 1) {
            os = android;
		} else if (userAgent.indexOf(linux) != -1) {
            os = linux;      
        } else if (userAgent.indexOf(windowsPhone) != - 1) {
            os = windowsPhone;
        } else if (userAgent.indexOf(windows) != - 1) {
            os = windows;
        } else if (userAgent.indexOf(mac) != - 1) {
            os = mac;
        } else if (userAgent.indexOf(blackberry) != - 1) {
            os = blackberry;
        } else if (userAgent.indexOf(bada) != - 1) {
            os = bada;
        } else {
            os = other;
        }
        //localStorage.setItem('os', os);
        return os;
    };
    this.isTester = function() {
        var tester = localStorage.getItem("tester");
        if (tester !== null) {
            return tester;
        }
        return 0;
    };
    this.getJqueryFriendlyEncId = function(encId) {
        return encId.substring(0, encId.length - 1);
    };
    this.formatMoney = function(value, addSymbol) {
        var noDecimal = value.toFixed(0); //Round up
        for (var i = 0; i < Math.floor((noDecimal.length - (1 + i)) / 3); i++)
        {
            noDecimal = noDecimal.substring(0, noDecimal.length - (4 * i + 3)) + ',' + noDecimal.substring(noDecimal.length - (4 * i + 3));
        }
        if(addSymbol) {
            noDecimal = '<span class="WebRupee">Rs.</span>' + noDecimal;
        }
        return noDecimal;
    };
    this.showMessage = function(msg, success) {
        var bg = "#56d619";
        if(success == false) {
            bg = "#FA5D5D;"
        }
        $('.top_message').css('background', bg).text(msg).slideDown();
        setTimeout(function() {
            $('.top_message').slideUp();
        }, 4000);
    };
};

//Count of the number of times viewed.
//Not to be confused with views in MVC context.
Views = function() {
    var baseKey = 'myViews';
    var syncId = 1;
    this.getBaseKey = function() {
        return baseKey;
    };

    this.getSyncUrl = function() {
        return '/view';
    };

    this.isTemp = function() {
        return true;
    };

    this.getAllViews = function() {
        var allViews = '';
        var storageKey = utilObj.getLocalKey(baseKey);
        var views = localStorage.getItem(storageKey);
        if (views !== null) {
            allViews += views;
        }

        storageKey = utilObj.getTempKey(baseKey);
        views = localStorage.getItem(storageKey);
        if (views !== null) {
            allViews += views;
        }

        return allViews.split(";");
    };

    this.addToViews = function(inspirationId) {
        if (inspirationId === undefined) {
            console.log("Inspiration id is not defined");
            return false;
        }
        var allViews = this.getAllViews();
        //Making the choice of not incrementing view count for the same inspiration id
        //for a person between sync intervals
        for (var i = 0; i < allViews.length; i++) {
            if (inspirationId == allViews[i]) {
                return false;
            }
        }

        var storageKey = utilObj.getLocalKey(baseKey);
        if (syncObj.isSyncInProgress(syncId)) {
            storageKey = utilObj.getTempKey(baseKey);
        }
        var curValue = localStorage.getItem(storageKey);
        if (curValue === null) {
            curValue = '';
        }
        curValue += inspirationId + ';';
        localStorage.setItem(storageKey, curValue);

        return true;
    };
};
