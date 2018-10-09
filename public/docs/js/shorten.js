function ShortenController(input,output,submitButton,copyButton){
	this.input = input;
	this.output = output;
	this.submitButton = submitButton;
	this.copyButton = copyButton;

	this.input.value = '';
	this.output.value = '';
	this.output.disabled = true;
	this.copyButton.disabled = true;
	this.submitButton.disabled = true;

	var that = this;
	this.input.addEventListener('input',function(e){
	  that.output.value = '';
	  that.output.disabled = true;
	  that.copyButton.disabled = true;
	  that.submitButton.disabled = e.target.value.length === 0;
	});
}

ShortenController.prototype.submit = function(){
	var that = this;
	fetch('https://small.ml/urlshortener/v1/url',{
		method: "POST",
		mode: "cors",
		credentials: "same-origin",
		body: JSON.stringify({
			longUrl: that.input.value,
		})
	})
	.then(res => res.json())
	.then(json => {
		that.output.value = json.shortUrl;
		that.copyButton.disabled = false;
		that.output.disabled = false;
	})
	.catch(function(err){
		console.log(JSON.stringify(err));
	});
};

ShortenController.prototype.copy = function(){
	if(!this.output.value){ return };
	if(!this.output.value.length){ return };

	this.output.focus();
	this.output.select();
	document.execCommand("copy");
}