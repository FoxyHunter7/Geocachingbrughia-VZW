<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('static_site_contents', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->string('property');
            $table->foreign('property')->references('property')->on('static_site_properties')->onUpdate('cascade')->onDelete('cascade');
            $table->string('lang_code');
            $table->foreign('lang_code')->references('code')->on('languages')->onUpdate('cascade')->onDelete('cascade');
            $table->text('content');
            $table->unique(['property','lang_code']);

            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('static_site_contents');
    }
};
