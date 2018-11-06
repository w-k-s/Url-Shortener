let service = {
	getShortenedUrl: async (longUrl)=>{
		const res = await fetch('/urlshortener/v1/url',{
			method: "POST",
			mode: "cors",
			credentials: "same-origin",
			body: JSON.stringify({
				longUrl,
			})
		});
		const json = await res.json();
		if(!res.ok){
			throw json;
		}
		return json;
	},
	getOriginalUrl: async (shortUrl) => {
		const res = await fetch(`/urlshortener/v1/url?shortUrl=${shortUrl}`,{
			method: "GET"
		});
		const json = await res.json();

		if(!res.ok){
			throw json;
		}

		return json;
	}
};

let controller = {
	init: ()=>{
		shortenView.init();
		originalView.init();
	},
	getShortenedUrl: async (longUrl)=>{
		try{
			const json = await service.getShortenedUrl(longUrl);
			const { shortUrl } = json;
			shortenView.showShortenedUrl(shortUrl);
		}catch(e){
			const {message} = e;
			shortenView.showError(message);
		}
	},
	getOriginalUrl: async (shortUrl)=>{
		try{
			const json = await service.getOriginalUrl(shortUrl);
			const {longUrl} = json;
			originalView.showOriginalUrl(longUrl);
		}catch(e){
			const {message} = e;
			originalView.showError(message);
		}
	},
	redirectToUrl: (url)=>{
		if(!url && !url.length){ return };
		window.location = url;
	}
};

let shortenView = {
	init: ()=>{
		this.longUrlInput = document.getElementById("inputLongUrl");
		this.shortUrlOutput = document.getElementById("outputShortUrl");
		this.longUrlInputError = document.getElementById("inputLongUrlError");
		this.submitLongUrl = document.getElementById("btnSubmitLongUrl");
		this.copyButton = document.getElementById("btnCopy");

		this.submitLongUrl.addEventListener('click',(e)=>{
			controller.getShortenedUrl(this.longUrlInput.value).then();
		});

		this.longUrlInput.addEventListener('input',(e)=>{
			this.shortUrlOutput.value = '';
			this.shortUrlOutput.disabled = true;
			this.copyButton.disabled = true;
			this.longUrlInputError.style.display = 'none';
			this.submitLongUrl.disabled = e.target.value.length === 0;
		});

		this.copyButton.addEventListener('click',(e)=>{
			const shortenedUrl = this.shortUrlOutput.value;
			if(!shortenedUrl || !shortenedUrl.length){ return };

			this.shortUrlOutput.focus();
			this.shortUrlOutput.select();
			document.execCommand("copy");
		});

		this.longUrlInput.value = '';
		this.shortUrlOutput.value = '';
		this.shortUrlOutput.disabled = true;
		this.longUrlInputError.style.display = 'none';
		this.copyButton.disabled = true;
		this.submitLongUrl.disabled = true;
	},
	showError: (error)=>{
		this.longUrlInput.value = '';
		this.longUrlInputError.style.display = 'block';
		this.longUrlInputError.children[0].innerHTML  = error;
	},
	showShortenedUrl: (shortenedUrl)=>{
		this.shortUrlOutput.value = shortenedUrl;
		this.shortUrlOutput.disabled = false;
		this.copyButton.disabled = false;
	}
};

let originalView = {
	init: ()=>{
		this.shortUrlInput = document.getElementById("inputShortUrl");
		this.originalUrlOutput = document.getElementById("outputOriginalUrl");
		this.shortUrlInputError = document.getElementById("inputShortUrlError");
		this.submitShortUrl = document.getElementById("btnSubmitShortUrl");
		this.goButton = document.getElementById("btnVisit");

		this.shortUrlInput.addEventListener('input',(e)=>{
		  this.originalUrlOutput.value = '';
		  this.originalUrlOutput.disabled = true;
		  this.goButton.disabled = true;
		  this.shortUrlInputError.style.display = 'none';
		  this.submitShortUrl.disabled = e.target.value.length === 0;
		});

		this.goButton.addEventListener('click',(e)=>{
			controller.redirectToUrl(this.originalUrlOutput.value);
		});

		this.submitShortUrl.addEventListener('click',(e)=>{
			controller.getOriginalUrl(this.shortUrlInput.value).then();
		});

		this.shortUrlInput.value = '';
		this.originalUrlOutput.value = '';
		this.originalUrlOutput.disabled = true;
		this.goButton.disabled = true;
		this.submitShortUrl.disabled = true;
	},
	showError: (error)=>{
		this.shortUrlInput.value = '';
		this.shortUrlInputError.style.display = 'block';
		this.shortUrlInputError.children[0].innerHTML  = error;
	},
	showOriginalUrl: (originalUrl)=>{
		this.originalUrlOutput.value = originalUrl;
		this.originalUrlOutput.disabled = false;
		this.goButton.disabled = false;
	}
};

function init(){
	controller.init();
}