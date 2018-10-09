function GetOriginalController(input,output,error,submitButton,goButton){
	this.input = input;
	this.output = output;
	this.error = error;
	this.submitButton = submitButton;
	this.goButton = goButton;

	this.input.value = '';
	this.output.value = '';
	this.output.disabled = true;
	this.goButton.disabled = true;
	this.submitButton.disabled = true;

	this.input.addEventListener('input',(e)=>{
	  this.output.value = '';
	  this.output.disabled = true;
	  this.goButton.disabled = true;
	  this.error.style.display = 'none';
	  this.submitButton.disabled = e.target.value.length === 0;
	});
}

GetOriginalController.prototype.submit = function(){
	fetch(`https://small.ml/urlshortener/v1/url?shortUrl=${this.input.value}`,{
		method: "GET"
	})
	.then(async res => {
		let json = await res.json();
		return {res, json};
	})
	.then(resp => {
		if(!resp.res.ok){
			throw resp.json;
		}
		this.output.value = resp.json.longUrl;
		this.goButton.disabled = false;
		this.output.disabled = false;
	})
	.catch((err)=>{
		this.input.value = '';
		this.error.style.display = 'block';
		this.error.children[0].innerHTML  = err.message;
	});
};

GetOriginalController.prototype.visitOriginal = function(){
	if(!this.output.value){ return };
	if(!this.output.value.length){ return };

	window.location = this.output.value;
}