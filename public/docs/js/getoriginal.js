function GetOriginalController(input,output,submitButton,goButton){
	this.input = input;
	this.output = output;
	this.submitButton = submitButton;
	this.goButton = goButton;

	this.input.value = '';
	this.output.value = '';
	this.output.disabled = true;
	this.goButton.disabled = true;
	this.submitButton.disabled = true;

	var that = this;
	this.input.addEventListener('input',function(e){
	  that.output.value = '';
	  that.output.disabled = true;
	  that.goButton.disabled = true;
	  that.submitButton.disabled = e.target.value.length === 0;
	});
}

GetOriginalController.prototype.submit = function(){
	var that = this;
	fetch(`https://small.ml/urlshortener/v1/url?shortUrl=${that.input.value}`,{
		method: "GET"
	})
	.then(res => res.json())
	.then(json => {
		that.output.value = json.longUrl;
		that.goButton.disabled = false;
		that.output.disabled = false;
	})
	.catch(function(err){
		console.log(JSON.stringify(err));
	});
};

GetOriginalController.prototype.visitOriginal = function(){
	if(!this.output.value){ return };
	if(!this.output.value.length){ return };

	window.location = this.output.value;
}