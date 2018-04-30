; (function () {

	'use strict';

	var isMobile = {
		Android: function () {
			return navigator.userAgent.match(/Android/i);
		},
		BlackBerry: function () {
			return navigator.userAgent.match(/BlackBerry/i);
		},
		iOS: function () {
			return navigator.userAgent.match(/iPhone|iPad|iPod/i);
		},
		Opera: function () {
			return navigator.userAgent.match(/Opera Mini/i);
		},
		Windows: function () {
			return navigator.userAgent.match(/IEMobile/i);
		},
		any: function () {
			return (isMobile.Android() || isMobile.BlackBerry() || isMobile.iOS() || isMobile.Opera() || isMobile.Windows());
		}
	};


	var fullHeight = function () {

		if (!isMobile.any()) {
			$('.js-fullheight').css('height', $(window).height());
			$(window).resize(function () {
				$('.js-fullheight').css('height', $(window).height());
			});
		}
	};

	// Parallax
	var parallax = function () {
		$(window).stellar();
	};

	var contentWayPoint = function () {
		var i = 0;
		$('.animate-box').waypoint(function (direction) {

			if (direction === 'down' && !$(this.element).hasClass('animated-fast')) {

				i++;

				$(this.element).addClass('item-animate');
				setTimeout(function () {

					$('body .animate-box.item-animate').each(function (k) {
						var el = $(this);
						setTimeout(function () {
							var effect = el.data('animate-effect');
							if (effect === 'fadeIn') {
								el.addClass('fadeIn animated-fast');
							} else if (effect === 'fadeInLeft') {
								el.addClass('fadeInLeft animated-fast');
							} else if (effect === 'fadeInRight') {
								el.addClass('fadeInRight animated-fast');
							} else {
								el.addClass('fadeInUp animated-fast');
							}

							el.removeClass('item-animate');
						}, k * 100, 'easeInOutExpo');
					});

				}, 50);

			}

		}, { offset: '85%' });
	};


	// sleep time expects milliseconds
	function sleep(time) {
		return new Promise((resolve) => setTimeout(resolve, time));
	}


	var loadData = function () {

		//TODO : Link API ( this is just a mocked up data )
		// var data_mocked = {
		// 	"user":
		// 	{
		// 		"id":"59ce3315087793ebafd92ea9",
		// 		"fullname":"Joscha Raue",
		// 		"image_url":"https://ca.slack-edge.com/T50JEDNKG-U4Z7D2BNC-8a4c441778a0-512"
		// 	},
		// 	"action":
		// 	{
		// 		"id":"59ce335b087793ebafd92ec0",
		// 		"description":"Hug everyone in the office"
		// 	},
		// 	"deadline":"2017-09-30T01:00:00+01:00"
		// }

		$.ajax({
			url: "/current_task",
			success: function (data) {
				var dateObject = new Date(data.deadline)
				dateObject = new Date(dateObject.toUTCString())
				$("#user-picture").css("background", "url('" + data.user.image_url + "')");
				$("#username").html(data.user.fullname);
				$("#action").html(data.action.message);
				startCountdown(dateObject);

			},
			error: function () {
				$("#action").html("<p> something went wrong :( </p> <p> please refresh the page to retry </p>");
				$("#username").css("display", "none");
				$("i").css("display", "none");
			}
		});
	}

	// Consume Api
	var getAction = function () {

		loadData();

		$(".fh5co-loader").fadeOut("slow");

	}

	// Loading page
	var loaderPage = function () {
		getAction();
	};

	//start countdown
	var startCountdown = function (countdown_start) {
		$("#main-example")
			.countdown(countdown_start, function (event) {
				$(this).text(
					event.strftime('%D days %H:%M:%S')
				);
			});
	}

	//Init countdown
	var initCountdown = function () {
		$("#main-example")
			.on('finish.countdown', function () {
				loadData()
			});

	}

	function getRandomInt(min, max) {
		return Math.floor(Math.random() * (max - min) + min);
	}

	//random background color
	var randomBackground = function () {
		var colorsPool = ['rgba(255, 50, 2, 0.8)', 'rgba(255, 50, 2, 0.8)', 'rgba(200, 50, 70, 0.8)',
			'rgba(70, 10, 20, 0.8)', 'rgba(200, 210, 20, 0.8)', 'rgba(80, 10, 20, 0.8)']
		var randomInt = getRandomInt(0, colorsPool.length)
		$(".overlay").css("background", colorsPool[randomInt])
	}




	$(function () {
		contentWayPoint();
		initCountdown();
		loaderPage();
		fullHeight();
		parallax();
		randomBackground();
	});



}());
