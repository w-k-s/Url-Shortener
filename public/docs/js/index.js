var shortenController;
var getOriginalController;

function init(){
	shortenController = new ShortenController(
		document.getElementById("inputLongUrl"),
		document.getElementById("outputShortUrl"),
		document.getElementById("btnSubmitLongUrl"),
		document.getElementById("btnCopy")
	)

	getOriginalViewModel = new GetOriginalController(
		document.getElementById("inputShortUrl"),
		document.getElementById("outputOriginalUrl"),
		document.getElementById("btnSubmitShortUrl"),
		document.getElementById("btnVisit")
	)
}

function shorten(){
	shortenController.submit();
}

function copyShortUrl(){
	shortenController.copy();
}

function getOriginal(){
	getOriginalController.submit();
}

function visitOriginal(){
	getOriginalController.visitOriginal();
}