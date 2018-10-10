function ShortenController(input,output,error,submitButton,copyButton){
	this.input = input;
	this.output = output;
	this.error = error;
	this.submitButton = submitButton;
	this.copyButton = copyButton;

	this.input.value = '';
	this.output.value = '';
	this.error.style.display = 'none';
	this.output.disabled = true;
	this.copyButton.disabled = true;
	this.submitButton.disabled = true;

	var that = this;
	this.input.addEventListener('input',function(e){
	  that.output.value = '';
	  that.output.disabled = true;
	  that.copyButton.disabled = true;
	  that.error.style.display = 'none';
	  that.submitButton.disabled = e.target.value.length === 0;
	});
}

ShortenController.prototype.submit = function(){
	var that = this;
	fetch('/urlshortener/v1/url',{
		method: "POST",
		mode: "cors",
		credentials: "same-origin",
		body: JSON.stringify({
			longUrl: that.input.value,
		})
	})
	.then(async res => {
		let json = await res.json();
		return {res, json};
	})
	.then(resp => {
		if(!resp.res.ok){
			throw resp.json;
		}

		that.output.value = resp.json.shortUrl;
		that.copyButton.disabled = false;
		that.output.disabled = false;
	})
	.catch(function(err){
		that.input.value = '';
		that.error.style.display = 'block';
		that.error.children[0].innerHTML  = err.message;
	});
};

ShortenController.prototype.copy = function(){
	if(!this.output.value){ return };
	if(!this.output.value.length){ return };

	this.output.focus();
	this.output.select();
	document.execCommand("copy");
}